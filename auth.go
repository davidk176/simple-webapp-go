/*
Enthält alle wichtigen Funtkionen für due User-Authentifizierung mit Google OAuth2.
Initialisiert die User-Session (gorilla sessions)
Liest und schreibt Cookies
*/

package main

import (
	"encoding/json"
	"github.com/davidk176/simple-webapp-go/utils"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/quasoft/memstore"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	tokenval "google.golang.org/api/oauth2/v2"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type Token struct {
	accesstoken  string
	refreshtoken string
	idtoken      string
	expiry       time.Time
}

var (
	googleOauthConfig *oauth2.Config
	httpClient        = &http.Client{}
	store             *memstore.MemStore
)

const (
	googleOAuthApiV2 = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="
	googleOAuthApiV4 = "https://www.googleapis.com/oauth2/v4/token"
)



/*
Initialisert die OAuthConfig und die Session
*/
func init() {
	ClientSecret, _ := utils.AccessSecretVersion("projects/test1-cc/secrets/CLIENT_SECRET/versions/latest")

	googleOauthConfig = &oauth2.Config{
		RedirectURL:  getRedirectUrl(),
		ClientID:     "345398956581-rq77v9k0l7uo0v7tvtgur21ld6ht3i8b.apps.googleusercontent.com",
		ClientSecret: *ClientSecret,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}

	//init Session mit zufälligen keys
	authKey := securecookie.GenerateRandomKey(64)
	encryptionKey := securecookie.GenerateRandomKey(32)

	store = memstore.NewMemStore(
		authKey,
		encryptionKey,
	)

	store.Options = &sessions.Options{
		MaxAge:   60 * 15, //Session läuft nach 15 min ab
		HttpOnly: true,    //sichert Cookie gegen Script-Zugriffe
		Secure:   true,    //erlaubt nur https
	}

}

/*
Ermittelt Redirect-Url zur Unterscheidung von localhost und App Engine
*/
func getRedirectUrl() (url string) {
	url = os.Getenv("OAUTH_REDIRECT_URL")
	if url == "" {
		return "http://localhost:8080/callback"
	}
	return url
}

/*
Wird von Google OAuth nach erfolgreichem Sign-In gerufen. Bekommt in Request den einmaligen Authentifizierungscode.
Dieser wird gegen den Token getauscht. Anschließend werden Userinformationen ermittelt.
*/
func handleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	//neue session anlegen
	session, err := store.Get(r, "session-name")
	writeDataInCache(w, r)
	if err != nil {
		log.Print(err)
	}
	state, _ := r.Cookie("oauthstate")
	log.Print(state)

	//vergleicht state aus cookie und state aus r zum Schutz vor XSRF
	if r.FormValue("state") != state.Value {
		log.Println("invalid oauth google state")
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return
	}

	//einmaligen auth code aus r holen
	code := r.FormValue("code")

	//neue Eingabe für Request zur Ermittlung des eigentlichen Tokens erstellen
	data := url.Values{}
	data.Set("code", code)                                    //der authorization_code aus dem ersten request
	data.Set("client_id", googleOauthConfig.ClientID)         //id des gcloud projekts
	data.Set("client_secret", googleOauthConfig.ClientSecret) //geheimes secret, in Google Secret Manager
	data.Set("redirect_uri", googleOauthConfig.RedirectURL)   //die redirect uri
	data.Set("grant_type", "authorization_code")              //grant über authorization_code
	data.Set("access_type", "offline")                        //offline für refresh_token
	log.Print(strings.NewReader(data.Encode()))

	//request an OAuth2, Refresh_Token in Session speichern
	jsonresult := callOAuthTokenUri(data)
	refresh_token := jsonresult["refresh_token"].(string)
	session.Values["refresh_token"] = refresh_token
	log.Print(refresh_token)

	//id_token in cookie setzen
	token := Token{
		accesstoken: jsonresult["access_token"].(string),
		idtoken:     jsonresult["id_token"].(string),
		expiry:      time.Now().Add(time.Duration(jsonresult["expires_in"].(float64)) * time.Second),
	}
	utils.GenerateTokenCookie(w, "idtoken", token.idtoken, token.expiry)

	//ermittelt User-Informationen von Google und speichert diese in Session
	responseuser, _ := http.Get(googleOAuthApiV2 + token.accesstoken)
	user, _ := ioutil.ReadAll(responseuser.Body)
	userStr := string(user)
	log.Print(userStr)
	var usermap map[string]interface{}
	err = json.Unmarshal([]byte(userStr), &usermap)
	log.Print(usermap)
	session.Values["username"] = usermap["name"]
	session.Values["picture"] = usermap["picture"]

	session.Save(r, w)
	log.Print("session saved")

	//Weiterleitung zu Shop
	log.Print(err)
	http.Redirect(w, r, "/shop", http.StatusFound)
	return
}

/*
Ausgangspunkt des Logins. Generiert State-Cookie und leitet an OAuth-Url weiter
*/
func loginHandler(w http.ResponseWriter, r *http.Request) {
	state := utils.GenerateStateCookie(w)
	url := googleOauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	log.Print("redirect to " + url)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

/*
Validiert den Id-Token des Requests. Falls kein Id-Token vorhanden, wird mit dem Refresh-Token aus der Session ein neuer geholt.
*/
func verifyIdToken(t string, w http.ResponseWriter, r *http.Request) bool {

	//open session
	session, err := store.Get(r, "session-name")
	refresh_token := session.Values["refresh_token"]

	//wenn kein refresh_token vorhanden, neuer Login notwendig
	if refresh_token == nil {
		log.Print("refresh_token is null --> login")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return false
	}
	log.Print("refresh_token from session " + refresh_token.(string))

	//leerer Token -> vermutlich Cookie expired -> hole neuen token
	if t == "" {
		log.Print("Cookie expired; --> refresh")
		data := url.Values{}
		data.Set("client_id", googleOauthConfig.ClientID)
		data.Set("client_secret", googleOauthConfig.ClientSecret)
		data.Set("refresh_token", refresh_token.(string))
		data.Set("grant_type", "refresh_token")
		result := callOAuthTokenUri(data)
		log.Print(result["access_token"])
		utils.GenerateTokenCookie(w, "idtoken", result["id_token"].(string), time.Now().Add(time.Duration(result["expires_in"].(float64))*time.Second))
		return true
	}
	//validiert Token mittels Aufruf von OAuth2-API
	oauth2Service, err := tokenval.New(httpClient)
	tokenInfoCall := oauth2Service.Tokeninfo()
	tokenInfoCall.IdToken(t)
	tokenInfo, err := tokenInfoCall.Do()
	log.Print(tokenInfo)
	log.Print(err)
	if err != nil {
		log.Print("Token invalid!")
		return false
	}
	return true
}

/*
Ruft googleapis/oauth2/v4/token mit den übergebenen Parametern.
*/
func callOAuthTokenUri(data url.Values) map[string]interface{} {
	client := &http.Client{}
	r, err := http.NewRequest("POST", googleOAuthApiV4, strings.NewReader(data.Encode())) // URL-encoded payload
	if err != nil {
		log.Fatal(err)
	}
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	res, err := client.Do(r)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(res.Status)
	defer res.Body.Close()

	//erstelle Map aus response
	resp, err := ioutil.ReadAll(res.Body)
	var result map[string]interface{}
	err = json.Unmarshal(resp, &result)
	log.Print(result)
	return result
}

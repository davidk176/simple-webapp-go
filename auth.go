package main

import (
	"encoding/base64"
	"encoding/json"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	tokenval "google.golang.org/api/oauth2/v2"
	"io/ioutil"
	"log"
	"math/rand"
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
	store             *sessions.CookieStore
)

const googleOAuthApi = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

func init() {
	ClientSecret, _ := accessSecretVersion("projects/test1-cc/secrets/CLIENT_SECRET/versions/latest")

	googleOauthConfig = &oauth2.Config{
		RedirectURL:  getRedirectUrl(),
		ClientID:     "345398956581-rq77v9k0l7uo0v7tvtgur21ld6ht3i8b.apps.googleusercontent.com",
		ClientSecret: *ClientSecret,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}

	//init session
	authKey := securecookie.GenerateRandomKey(64)
	encryptionKey := securecookie.GenerateRandomKey(32)

	store = sessions.NewCookieStore(
		authKey,
		encryptionKey,
	)

	store.Options = &sessions.Options{
		MaxAge:   60 * 15,
		HttpOnly: true,
	}
}

func getRedirectUrl() (url string) {
	url = os.Getenv("OAUTH_REDIRECT_URL")
	if url == "" {
		return "http://localhost:8080/callback"
	}
	return url
}

func handleGoogleCallback(w http.ResponseWriter, r *http.Request) {

	//neue session anlegen
	session, err := store.Get(r, "session-name")
	if err != nil {
		log.Print(err)
	}

	//state aus Cookie holen und mit state aus request vergleichen
	state, _ := r.Cookie("oauthstate")
	log.Print(state)

	if r.FormValue("state") != state.Value {
		log.Println("invalid oauth google state")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	//einmaligen auth code aus r holen
	code := r.FormValue("code")

	//neue Eingabe für Request für id_token erstellen
	data := url.Values{}
	data.Set("code", code)
	data.Set("client_id", googleOauthConfig.ClientID)
	data.Set("client_secret", googleOauthConfig.ClientSecret)
	data.Set("redirect_uri", googleOauthConfig.RedirectURL)
	data.Set("grant_type", "authorization_code")
	data.Set("access_type", "offline")
	//	data.Set("approval_prompt", "force")
	//	data.Set("prompt", "select_account") //select_account consent
	log.Print(strings.NewReader(data.Encode()))

	//request oauth2
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
	generateTokenCookie(w, "idtoken", token.idtoken, token.expiry)

	responseuser, _ := http.Get(googleOAuthApi + token.accesstoken)
	user, _ := ioutil.ReadAll(responseuser.Body)
	userStr := string(user)

	log.Print(userStr)
	session.Save(r, w)
	http.Redirect(w, r, "/shop", http.StatusPermanentRedirect)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	state := generateStateCookie(w)
	url := googleOauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	log.Print(url)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func generateStateCookie(w http.ResponseWriter) string {
	var exp = time.Now().Add(365 * 24 * time.Hour)
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: exp, HttpOnly: true}
	http.SetCookie(w, &cookie)
	return state
}

func generateTokenCookie(w http.ResponseWriter, n string, value string, e time.Time) {
	cookie := http.Cookie{Name: n, Value: value, Expires: e, HttpOnly: true}
	http.SetCookie(w, &cookie)
}

func verifyIdToken(t string, w http.ResponseWriter, r *http.Request) bool {

	//open session
	session, err := store.Get(r, "session-name")
	refresh_token := session.Values["refresh_token"].(string)
	log.Print("refresh_token from session " + refresh_token)

	if t == "" {
		log.Print("Cookie expired; --> refresh")
		data := url.Values{}
		data.Set("client_id", googleOauthConfig.ClientID)
		data.Set("client_secret", googleOauthConfig.ClientSecret)
		data.Set("refresh_token", refresh_token)
		data.Set("grant_type", "refresh_token")
		result := callOAuthTokenUri(data)
		log.Print(refresh_token)
		log.Print(result["access_token"])
		return true
	}

	oauth2Service, err := tokenval.New(httpClient)
	tokenInfoCall := oauth2Service.Tokeninfo()
	tokenInfoCall.IdToken(t)
	tokenInfo, err := tokenInfoCall.Do()
	log.Print(tokenInfo)
	log.Print(err)
	if err != nil {
		return false
	}
	return true
}

func getInfoFromCookie(c *http.Cookie) string {
	if c != nil {
		log.Print("Token from Cookie: " + c.Value)
		return c.Value
	}
	log.Print("Cookie is null")
	return ""
}

func refresh(w http.ResponseWriter, r *http.Request) (t Token) {
	return t
}

func callOAuthTokenUri(data url.Values) map[string]interface{} {
	client := &http.Client{}
	r, err := http.NewRequest("POST", "https://www.googleapis.com/oauth2/v4/token", strings.NewReader(data.Encode())) // URL-encoded payload
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

	//erstelle map aus response
	resp, err := ioutil.ReadAll(res.Body)
	var jsonresult map[string]interface{}
	err = json.Unmarshal(resp, &jsonresult)
	log.Print(jsonresult)
	return jsonresult
}

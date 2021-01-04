package main

import (
	"context"
	"encoding/base64"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type Token struct {
	accesstoken  string
	refreshtoken string
	expiry       time.Time
}

var (
	googleOauthConfig *oauth2.Config
	httpClient        = &http.Client{}
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
}

func getRedirectUrl() (url string) {
	url = os.Getenv("OAUTH_REDIRECT_URL")
	if url == "" {
		return "http://localhost:8080/callback"
	}
	return url
}

func handleGoogleCallback(w http.ResponseWriter, r *http.Request) {

	state, err := r.Cookie("oauthstate")
	log.Print(state)

	if r.FormValue("state") != state.Value {
		log.Println("invalid oauth google state")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	res, err := http.Get(googleOAuthApi + token.AccessToken)

	if err != nil {
		return
	}
	log.Print(token.AccessToken)
	log.Print(token.Expiry.String())
	log.Print(token.RefreshToken)
	T := Token{
		accesstoken:  token.AccessToken,
		refreshtoken: token.RefreshToken,
		expiry:       token.Expiry,
	}

	generateTokenCookie(w, "accesstoken", T)
	user, _ := ioutil.ReadAll(res.Body)
	userStr := string(user)

	log.Print(userStr)
	http.Redirect(w, r, "/shop?auth="+token.AccessToken, http.StatusPermanentRedirect)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	state := generateStateCookie(w)
	url := googleOauthConfig.AuthCodeURL(state)
	log.Print(url)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func generateStateCookie(w http.ResponseWriter) string {
	var exp = time.Now().Add(365 * 24 * time.Hour)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: exp}
	http.SetCookie(w, &cookie)
	return state
}

func generateTokenCookie(w http.ResponseWriter, n string, t Token) {
	cookie := http.Cookie{
		Name:       n,
		Value:      t.accesstoken,
		Path:       "",
		Domain:     "",
		Expires:    t.expiry,
		RawExpires: "",
		MaxAge:     0,
		Secure:     false,
		HttpOnly:   false,
		SameSite:   0,
		Raw:        "",
		Unparsed:   nil,
	}
	http.SetCookie(w, &cookie)
}

func verifyIdToken(t string) bool {
	//TODO: Implement Verification
	if t == "" {
		return false
	}
	return true
}

package main

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"log"
	"net/http"
	"os"
)

var (
	googleOauthConfig *oauth2.Config
	oauthStateString  = "pseudo-random"
)

func init() {
	ClientSecret, _ := accessSecretVersion("projects/test1-cc/secrets/DB_SQL_PW/versions/latest")

	googleOauthConfig = &oauth2.Config{
		RedirectURL:  getRedirectUrl(),
		ClientID:     "345398956581-rq77v9k0l7uo0v7tvtgur21ld6ht3i8b.apps.googleusercontent.com",
		ClientSecret: *ClientSecret,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
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

	state := r.FormValue("state")
	log.Print(state)

	code := r.FormValue("code")

	token, err := googleOauthConfig.Exchange(oauth2.NoContext, code)

	if err != nil {
		return
	}
	log.Print(token.AccessToken)
	log.Print(token.Expiry.String())
	log.Print(token.RefreshToken)

	http.Redirect(w, r, "/shop", http.StatusPermanentRedirect)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	log.Print(url)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

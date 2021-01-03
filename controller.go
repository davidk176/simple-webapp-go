package main

import (
	"golang.org/x/oauth2/google"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

	"golang.org/x/oauth2"
	_ "golang.org/x/oauth2/google"
)

type PageVar struct {
	Title    string
	Response string
	Artikel  []Artikel
}

type Artikel struct {
	Name string
	Anz  int64
}

var (
	googleOauthConfig *oauth2.Config
	oauthStateString  = "pseudo-random"
)

func init() {
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  getRedirectUrl(),
		ClientID:     "345398956581-rq77v9k0l7uo0v7tvtgur21ld6ht3i8b.apps.googleusercontent.com",
		ClientSecret: "hi593_XKTINSKQuZQ741K1MK",
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

func shoppingHandler(w http.ResponseWriter, _ *http.Request) {
	log.Print("Start shoppingHandler")

	Var := PageVar{
		Title: "MyShop",
	}
	t, err := template.ParseFiles("templates/shop1.html")

	if err != nil {
		log.Print("Error parsing template: ", err)
	}
	err = t.Execute(w, Var)
}

func artikelHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("Start artikelHandler")
	_ = r.ParseForm()

	log.Print(r)
	name := r.Form.Get("name")
	menge, _ := strconv.ParseInt(r.Form.Get("menge"), 10, 64)

	Var := PageVar{
		Title:    "MyShop",
		Response: name,
	}

	Artikel := Artikel{
		Name: name,
		Anz:  menge,
	}
	log.Print(Artikel)
	addArtikelToDatabase(Artikel)

	t, err := template.ParseFiles("templates/shop1.html")

	if err != nil {
		log.Print("Error parsing template: ", err)
	}
	Var.Artikel = getArtikelFromDatabase()
	err = t.Execute(w, Var)

}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/start.html")

	if err != nil {
		log.Print("Error parsing template: ", err)
	}

	err = t.Execute(w, nil)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	log.Print(url)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
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

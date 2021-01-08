package main

import (
	_ "golang.org/x/oauth2/google"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type PageVar struct {
	Title    string
	Response string
	Name     string
	Picture  string
	Username string
	Artikel  []Artikel
}

type Artikel struct {
	Id   int64
	Name string
	Anz  int64
}

func shoppingHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("Start shoppingHandler")
	session, _ := store.Get(r, "session-name")

	cookie, _ := r.Cookie("idtoken")
	cv := getInfoFromCookie(cookie)

	if !verifyIdToken(cv, w, r) {
		http.Redirect(w, r, "/error", http.StatusPermanentRedirect)
	}

	pv := PageVar{
		Title:    "MyShop",
		Picture:  session.Values["picture"].(string),
		Username: session.Values["username"].(string),
	}
	t, err := template.ParseFiles("templates/shop1.html")

	if err != nil {
		log.Print("Error parsing template: ", err)
	}
	err = t.Execute(w, pv)
}

func artikelHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("Start artikelHandler")
	session, _ := store.Get(r, "session-name")

	cookie, _ := r.Cookie("idtoken")
	log.Print("Token from Cookie: " + cookie.Value)
	cv := getInfoFromCookie(cookie)
	if !verifyIdToken(cv, w, r) {
		return
	}

	_ = r.ParseForm()
	log.Print(r)
	name := r.Form.Get("name")
	menge, _ := strconv.ParseInt(r.Form.Get("menge"), 10, 64)

	pv := PageVar{
		Title:    "MyShop",
		Response: name,
		Picture:  session.Values["picture"].(string),
		Username: session.Values["username"].(string),
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
	pv.Artikel = getArtikelFromDatabase()
	err = t.Execute(w, pv)

}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/start.html")

	if err != nil {
		log.Print("Error parsing template: ", err)
	}

	err = t.Execute(w, nil)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("Start deleteHandler")
	session, _ := store.Get(r, "session-name")

	pv := PageVar{
		Title:    "MyShop",
		Picture:  session.Values["picture"].(string),
		Username: session.Values["username"].(string),
	}

	cookie, _ := r.Cookie("idtoken")
	cv := getInfoFromCookie(cookie)
	if !verifyIdToken(cv, w, r) {
		return
	}

	_ = r.ParseForm()
	log.Print(r)
	id := r.Form.Get("deleteid")
	log.Print("delete id: " + id)
	deleteArtikelFromDatabase(id)

	t, err := template.ParseFiles("templates/shop1.html")

	if err != nil {
		log.Print("Error parsing template: ", err)
	}
	pv.Artikel = getArtikelFromDatabase()
	err = t.Execute(w, pv)

}

func errorHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/err.html")

	if err != nil {
		log.Print("Error parsing template: ", err)
	}
	err = t.Execute(w, nil)
}

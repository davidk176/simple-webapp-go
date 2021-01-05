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
	Artikel  []Artikel
}

type Artikel struct {
	Id   int64
	Name string
	Anz  int64
}

func shoppingHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("Start shoppingHandler")

	cookie, _ := r.Cookie("accesstoken")
	log.Print("Token from Cookie: " + cookie.Value)

	if !verifyIdToken(cookie.Value) {
		return
	}

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

	cookie, _ := r.Cookie("accesstoken")
	log.Print("Token from Cookie: " + cookie.Value)
	if !verifyIdToken(cookie.Value) {
		return
	}

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

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("Start deleteHandler")

	cookie, _ := r.Cookie("accesstoken")
	log.Print("Token from Cookie: " + cookie.Value)
	if !verifyIdToken(cookie.Value) {
		return
	}

}

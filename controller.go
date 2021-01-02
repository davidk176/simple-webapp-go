package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
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

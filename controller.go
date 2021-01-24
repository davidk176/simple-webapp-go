package main

import (
	"github.com/davidk176/simple-webapp-go/utils"
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
	session, err := store.Get(r, "session-name")
	log.Print("Session", session)

	if err != nil {
		log.Print(err)
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
		return
	}

	cookie, _ := r.Cookie("idtoken")
	cv := utils.GetInfoFromCookie(cookie)

	if !verifyIdToken(cv, w, r) {
		http.Redirect(w, r, "/error", http.StatusPermanentRedirect)
		return
	}

	pv := PageVar{
		Title:    "MyShop",
		Picture:  session.Values["picture"].(string),
		Username: session.Values["username"].(string),
	}
	pv.Artikel = getArtikelFromDatabase()
	t, err := template.ParseFiles("templates/shop1.html")

	if err != nil {
		log.Print("Error parsing template: ", err)
	}
	err = t.Execute(w, pv)
}

func artikelHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("Start artikelHandler")
	session, err := store.Get(r, "session-name")

	if err != nil {
		log.Print(err)
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
		return
	}

	cookie, _ := r.Cookie("idtoken")
	cv := utils.GetInfoFromCookie(cookie)
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
	session, err := store.Get(r, "session-name")

	cookie, _ := r.Cookie("idtoken")
	cv := utils.GetInfoFromCookie(cookie)
	if !verifyIdToken(cv, w, r) {
		return
	}

	if err != nil {
		log.Print(err)
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
		return
	}
	cookie, _ := r.Cookie("idtoken")
	cv := utils.GetInfoFromCookie(cookie)
	if !verifyIdToken(cv, w, r) {
		return
	}

	pv := PageVar{
		Title:    "MyShop",
		Picture:  session.Values["picture"].(string),
		Username: session.Values["username"].(string),
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

/*func errorHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/default_error.html")

	if err != nil {
		log.Print("Error parsing template: ", err)
	}
	err = t.Execute(w, nil)
}*/

package main

import (
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
)

type PageVar struct {
	Title    string
	Response string
}

type Artikel struct {
	Name string
	Anz  int64
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Print("Using  Port " + port)

	mux := http.NewServeMux()

	//mux.HandleFunc("/", greetingHandler)
	mux.HandleFunc("/", shoppingHandler)
	mux.HandleFunc("/add", articelHandler)
	_ = http.ListenAndServe(":"+port, mux)
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

func articelHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("Start articelHandler")
	_ = r.ParseForm()

	log.Print(r)
	name := r.Form.Get("name")
	menge, _ := strconv.ParseInt(r.Form.Get("menge"), 10, 64)

	Var := PageVar{
		Title:    "MyShop",
		Response: name,
	}

	t, err := template.ParseFiles("templates/shop1.html")

	if err != nil {
		log.Print("Error parsing template: ", err)
	}
	err = t.Execute(w, Var)

	artikel := Artikel{
		Name: name,
		Anz:  menge,
	}
	log.Print(artikel)
	writeToDatabase(artikel)

}

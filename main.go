package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

type Greeting struct {
	Name string
	Time string
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Print("Using Port " + port)

	mux := http.NewServeMux()

	mux.HandleFunc("/", greetingHandler)
	http.ListenAndServe(":"+port, mux)
}

func greetingHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("Start greetingHandler")
	greeting := Greeting{
		Name: "any",
		Time: time.Now().Format(time.Stamp),
	}

	t, err := template.ParseFiles("templates/start.html")
	writeToDatabase("Buch", 2)

	if err != nil {
		log.Print("Error parsing template: ", err)
	}
	err = t.Execute(w, greeting)
}

func writeToDatabase(name string, anzahl int) {

	log.Print("Write to DB")
	db, err := sql.Open("mysql", "admin:admin@tcp(127.0.0.1:3306)/webapp")
	if err != nil {
		log.Print(err.Error())
	}

	insert, err := db.Query("INSERT INTO Artikel (Name, Anzahl) VALUES ('test', 2)")

	if insert != nil {
		log.Print(insert.Columns())
	}
	if err != nil {
		log.Fatal(err.Error())
	}

	defer db.Close()

}

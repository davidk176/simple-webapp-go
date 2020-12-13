package main

import (
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

	mux := http.NewServeMux()

	mux.HandleFunc("/", greetingHandler)
	http.ListenAndServe(":"+port, mux)
}

func greetingHandler(w http.ResponseWriter, r *http.Request) {
	greeting := Greeting{
		Name: "any",
		Time: time.Now().Format(time.Stamp),
	}

	t, err := template.ParseFiles("templates/start.html")

	if err != nil {
		log.Print("Error parsing template: ", err)
	}
	err = t.Execute(w, greeting)
}

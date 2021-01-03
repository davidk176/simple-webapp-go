package main

import (
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"os"
)

func main() {
	port := getPort()

	mux := http.NewServeMux()
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/login", loginHandler)
	mux.HandleFunc("/callback", handleGoogleCallback)

	mux.HandleFunc("/shop", shoppingHandler)
	mux.HandleFunc("/add", artikelHandler)

	_ = http.ListenAndServe(":"+port, mux)
}

func getPort() (port string) {
	port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Print("Using  Port " + port)
	return port
}

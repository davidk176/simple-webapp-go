package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/context"
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
	mux.HandleFunc("/delete", deleteHandler)

	mux.HandleFunc("/error", errorHandler)

	_ = http.ListenAndServe(":"+port, context.ClearHandler(mux))

}

func getPort() (port string) {
	port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Print("Using  Port " + port)
	return port
}

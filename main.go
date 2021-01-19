package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"google.golang.org/appengine"
	"log"
	"net/http"
	"os"
)

func init() {
	router := mux.NewRouter()
	router.HandleFunc("/", homeHandler)
	router.HandleFunc("/login", loginHandler)
	router.HandleFunc("/callback", handleGoogleCallback)

	router.HandleFunc("/shop", shoppingHandler)
	router.HandleFunc("/add", artikelHandler)
	router.HandleFunc("/delete", deleteHandler)

	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(":"+getPort(), context.ClearHandler(router)))
}

func main() {
	appengine.Main()

}

func getPort() (port string) {
	port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Print("Using  Port " + port)
	return port
}

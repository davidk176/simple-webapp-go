package main

import (
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

func init() {

}

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/", homeHandler)
	router.HandleFunc("/login", loginHandler)
	router.HandleFunc("/callback", handleGoogleCallback)

	router.HandleFunc("/shop", shoppingHandler)
	router.HandleFunc("/add", artikelHandler)
	router.HandleFunc("/delete", deleteHandler)

	router.HandleFunc("/addCalculator", calculatorHandler)

	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(":"+getPort(), context.ClearHandler(router)))

}

func getPort() (port string) {
	port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Print("Using  Port " + port)
	return port
}

package main

import (
	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/appengine"
	"log"
	"net/http"
	"os"
)

func init() {

}

func main() {

	//	router := mux.NewRouter()
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/callback", handleGoogleCallback)

	http.HandleFunc("/shop", shoppingHandler)
	http.HandleFunc("/add", artikelHandler)
	http.HandleFunc("/delete", deleteHandler)

	http.ListenAndServe(":"+getPort(), nil)
	//http.Handle("/", router)
	//log.Fatal(http.ListenAndServe(":"+getPort(), context.ClearHandler(router)))
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

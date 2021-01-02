package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Greeting struct {
	Name string
	Time string
}

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
	http.ListenAndServe(":"+port, mux)
}

func greetingHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("Start greetingHandler")
	greeting := Greeting{
		Name: "any",
		Time: time.Now().Format(time.Stamp),
	}

	t, err := template.ParseFiles("templates/start.html")
	//writeToDatabase("Buch", 2)

	if err != nil {
		log.Print("Error parsing template: ", err)
	}
	err = t.Execute(w, greeting)
}

func shoppingHandler(w http.ResponseWriter, r *http.Request) {
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
	r.ParseForm()

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

func writeToDatabase(artikel Artikel) {

	log.Print("Write to DB: " + artikel.Name)
	//db, err := sql.Open("mysql", "admin:admin@tcp(127.0.0.1:3306)/webapp")
	db, err := initSocketConnectionPool()
	if err != nil {
		log.Print(err.Error())
	}

	insert, err := db.Query("INSERT INTO Artikel (Name, Anzahl) VALUES (?, ?)", artikel.Name, artikel.Anz)

	if insert != nil {
		log.Print(insert.Columns())
	}
	if err != nil {
		log.Fatal(err.Error())
	}

	defer db.Close()

}

func initSocketConnectionPool() (*sql.DB, error) {
	// [START cloud_sql_mysql_databasesql_create_socket]
	var (
		dbUser                 = os.Getenv("DB_USER")
		dbPwd                  = os.Getenv("DB_PASS")
		instanceConnectionName = os.Getenv("INSTANCE_CONNECTION_NAME")
		dbName                 = os.Getenv("DB_NAME")
	)

	socketDir, isSet := os.LookupEnv("DB_SOCKET_DIR")
	if !isSet {
		socketDir = "/cloudsql"
	}

	var dbURI string
	dbURI = fmt.Sprintf("%s:%s@unix(/%s/%s)/%s?parseTime=true", dbUser, dbPwd, socketDir, instanceConnectionName, dbName)

	// dbPool is the pool of database connections.
	dbPool, err := sql.Open("mysql", dbURI)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %v", err)
	}

	return dbPool, nil
	// [END cloud_sql_mysql_databasesql_create_socket]
}

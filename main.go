package main

import (
	"database/sql"
	"fmt"
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
	//db, err := sql.Open("mysql", "admin:admin@tcp(127.0.0.1:3306)/webapp")
	db, err := initSocketConnectionPool()
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

func initSocketConnectionPool() (*sql.DB, error) {
	// [START cloud_sql_mysql_databasesql_create_socket]
	var (
		dbUser                 = os.Getenv("DB_USER")                  // e.g. 'my-db-user'
		dbPwd                  = os.Getenv("DB_PASS")                  // e.g. 'my-db-password'
		instanceConnectionName = os.Getenv("INSTANCE_CONNECTION_NAME") // e.g. 'project:region:instance'
		dbName                 = os.Getenv("DB_NAME")                  // e.g. 'my-database'
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

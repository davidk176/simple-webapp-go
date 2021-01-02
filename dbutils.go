package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

func addArtikelToDatabase(artikel Artikel) {

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

func getArtikelFromDatabase() (artikel []Artikel) {
	log.Print("Get Artikel from DB")
	db, err := initSocketConnectionPool()
	if err != nil {
		log.Print(err.Error())
	}

	rows, _ := db.Query("SELECT * FROM Artikel")

	for rows.Next() {
		a := Artikel{}
		rows.Scan(&a.Name, &a.Anz, nil)
		artikel = append(artikel, a)
	}

	return artikel
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

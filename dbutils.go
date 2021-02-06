/*
Verbindung zur MySQL Datenbank.
*/

package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/davidk176/simple-webapp-go/utils"
	_ "github.com/go-sql-driver/mysql"

	"context"
)

func addArtikelToDatabase(artikel Artikel) {

	log.Print("Write to DB: " + artikel.Name)
	db, err := initSocketConnectionPool()
	log.Print("initialized DB socket ")
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
		rows.Scan(&a.Name, &a.Anz, &a.Id)
		artikel = append(artikel, a)
	}

	return artikel
}

func deleteArtikelFromDatabase(id string) (artikel Artikel) {
	log.Print("delete Artikel with id " + id)
	db, err := initSocketConnectionPool()
	if err != nil {
		log.Print(err.Error())
	}
	rows, _ := db.Query("DELETE FROM Artikel WHERE ID=?", id)

	rows.Next()
	a := Artikel{}
	rows.Scan(&a.Name, &a.Anz, &a.Id)

	//TEST
	fmt.Println("TestDelete")
	ctx := context.Background()
	createClient(ctx)

	fmt.Println(createClient(ctx))

	return artikel
}

func createClient(ctx context.Context) *firestore.Client {
	// Sets your Google Cloud Platform project ID.
	projectID := "webapp-shop-303617"

	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	fmt.Println(client)

	_, _, err2 := client.Collection("users").Add(ctx, map[string]interface{}{
		"first": "Ada",
		"last":  "Lovelace",
		"born":  1815,
	})
	if err2 != nil {
		log.Fatalf("Failed adding alovelace: %v", err2)
	}

	// Close client when done with
	defer client.Close()
	return client
}

func initSocketConnectionPool() (*sql.DB, error) {
	// [START cloud_sql_mysql_databasesql_create_socket]
	var (
		dbUser                 = os.Getenv("DB_USER")
		dbPwd, err             = utils.AccessSecretVersion("projects/webapp-shop-303617/secrets/DB_SQL_PW/versions/latest")
		instanceConnectionName = os.Getenv("INSTANCE_CONNECTION_NAME")
		dbName                 = os.Getenv("DB_NAME")
	)

	if err != nil {
		log.Print(err)
	}

	socketDir, isSet := os.LookupEnv("DB_SOCKET_DIR")
	if !isSet {
		socketDir = "/cloudsql"
	}

	var dbURI string
	dbURI = fmt.Sprintf("%s:%s@unix(/%s/%s)/%s?parseTime=true", dbUser, *dbPwd, socketDir, instanceConnectionName, dbName)

	// dbPool is the pool of database connections.
	dbPool, err := sql.Open("mysql", dbURI)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %v", err)
	}

	return dbPool, nil
	// [END cloud_sql_mysql_databasesql_create_socket]
}

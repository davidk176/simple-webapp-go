package main

import (
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"context"
	"database/sql"
	"fmt"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
	"log"
	"os"
)

func addArtikelToDatabase(artikel Artikel) {

	log.Print("Write to DB: " + artikel.Name)
	//db, err := sql.Open("mysql", "admin:admin@tcp(127.0.0.1:3306)/webapp")
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
		rows.Scan(&a.Name, &a.Anz, nil)
		artikel = append(artikel, a)
	}

	return artikel
}

func initSocketConnectionPool() (*sql.DB, error) {
	// [START cloud_sql_mysql_databasesql_create_socket]
	var (
		dbUser                 = os.Getenv("DB_USER")
		dbPwd, err             = accessSecretVersion("projects/test1-cc/secrets/pw-sqldb/versions/latest")
		instanceConnectionName = os.Getenv("INSTANCE_CONNECTION_NAME")
		dbName                 = os.Getenv("DB_NAME")
	)
	log.Print(dbPwd)

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

// accessSecretVersion accesses the payload for the given secret version if one
// exists. The version can be a version number as a string (e.g. "5") or an
// alias (e.g. "latest").
func accessSecretVersion(name string) (*string, error) {
	// name := "projects/my-project/secrets/my-secret/versions/5"
	// name := "projects/my-project/secrets/my-secret/versions/latest"
	log.Print("access Secret " + name)
	var secret string

	// Create the client.
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create secretmanager client: %v", err)
	}

	// Build the request.
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	}

	// Call the API.
	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to access secret version: %v", err)
	}
	secret = string(result.Payload.Data)

	log.Print(secret)
	return &secret, nil
}

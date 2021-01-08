/*
Zugriff auf Secrets aus Google Secret Manager
*/

package main

import (
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"context"
	"fmt"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
	"log"
)

//Methode mit leichten Änderungen aus Google Doku übernommen: https://cloud.google.com/secret-manager/docs/samples/secretmanager-access-secret-version#secretmanager_access_secret_version-go
// accessSecretVersion accesses the payload for the given secret version if one
// exists. The version can be a version number as a string (e.g. "5") or an projects/345398956581/secrets/pw-sqldb/versions/1 projects/345398956581/secrets/pw-sqldb
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

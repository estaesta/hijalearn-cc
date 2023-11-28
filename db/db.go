package db

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
)

const projectID = "festive-antenna-402105"

// InitFirestoreClient initializes a Firestore client.
func CreateClient(ctx context.Context) *firestore.Client {
	// Sets your Google Cloud Platform project ID.
	// projectID := os.Getenv("PROJECT_ID")

	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	// Close client when done with
	// defer client.Close()
	return client
}


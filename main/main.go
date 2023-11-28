package main

import (
	"context"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/estaesta/hijalearn/auth"
	"github.com/estaesta/hijalearn/db"
	"github.com/estaesta/hijalearn/handlers"
	"github.com/labstack/echo/v4"
)

var dbClient *firestore.Client

func main() {
	e := echo.New()

	// initialize firebase service and middleware
	// projectID := os.Getenv("PROJECT_ID")
	projectID := "festive-antenna-402105"
	firebaseService := auth.NewFirebaseService(projectID)
	firebaseMiddleware := auth.FirebaseMiddleware(firebaseService)

	// initialize firestore client
	dbClient = db.CreateClient(context.Background())
	defer dbClient.Close()

	// routes
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	// get user's learning progress
	getProgressUser := func(c echo.Context) error {
		return handlers.GetProgressUser(c, dbClient)
	}
	e.GET("/progress/", getProgressUser, firebaseMiddleware)

	// update user's learning progress
	updateProgressUser := func(c echo.Context) error {
		return handlers.UpdateProgressUser(c, dbClient)
	}
	e.PUT("/progress/", updateProgressUser, firebaseMiddleware)

	// initialize user's learning progress
	initProgressUser := func(c echo.Context) error {
		return handlers.InitProgressUser(c, dbClient)
	}
	e.POST("/progress/", initProgressUser, firebaseMiddleware)

	e.Logger.Fatal(e.Start(":8080"))
}

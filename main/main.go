package main

import (
	"context"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/estaesta/hijalearn/auth"
	"github.com/estaesta/hijalearn/db"
	"github.com/estaesta/hijalearn/handlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var dbClient *firestore.Client

func main() {
	e := echo.New()

	// initialize firebase service and middleware
	// projectID := os.Getenv("PROJECT_ID")
	projectID := "festive-antenna-402105"
	firebaseService := auth.NewFirebaseService(projectID)
	firebaseMiddleware := auth.FirebaseMiddleware(firebaseService)
	mlEndpoint := "https://hijalearn-ml-e6mqsjvzxq-et.a.run.app/predict"
	// mlEndpoint := "http://localhost:5000/predict"

	// initialize firestore client
	dbClient = db.CreateClient(context.Background())
	defer dbClient.Close()

	e.Use(middleware.Logger())

	// routes
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "API is running")
	})

	// get user's learning progress
	getProgressUser := func(c echo.Context) error {
		return handlers.GetProgressUser(c, dbClient)
	}
	e.GET("/api/v1/progress", getProgressUser, firebaseMiddleware)

	// get single module user's learning progress
	getProgressUserModule := func(c echo.Context) error {
		return handlers.GetProgressUserModule(c, dbClient)
	}
	e.GET("/api/v1/progress/:moduleId", getProgressUserModule, firebaseMiddleware)

	// update user's learning progress
	// updateProgressUser := func(c echo.Context) error {
	// 	return handlers.UpdateProgressUser(c, dbClient)
	// }
	// e.PUT("/api/v1/progress", updateProgressUser, firebaseMiddleware)

	// initialize user's learning progress
	initProgressUser := func(c echo.Context) error {
		return handlers.InitProgressUser(c, dbClient)
	}
	e.POST("/api/v1/progress", initProgressUser, firebaseMiddleware)

	// register
	register := func(c echo.Context) error {
		return handlers.Register(c, firebaseService, dbClient)
	}
	e.POST("/api/v1/register", register)

	// prediction
	predict := func(c echo.Context) error {
		return handlers.Predict(c, dbClient, mlEndpoint)
	}
	e.POST("/api/v1/prediction", predict, firebaseMiddleware)

	//Update Profile
	updateProfile := func(c echo.Context) error {
		return handlers.UpdateProfile(c, firebaseService)
	}
	e.POST("/api/v1/update_profile", updateProfile, firebaseMiddleware)

	e.Logger.Fatal(e.Start(":8080"))
}

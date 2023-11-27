package auth

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/labstack/echo/v4"
)

type FirebaseService struct {
	firebase *firebase.App
	app      *auth.Client
}

func NewFirebaseService(projectID string) *FirebaseService {
	ctx := context.Background()
	config := &firebase.Config{
		ProjectID: projectID,
	}

	firebase, err := firebase.NewApp(ctx, config)
	if err != nil {
		log.Fatalln(err)
	}

	app, err := firebase.Auth(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	return &FirebaseService{
		firebase: firebase,
		app:      app,
	}
}

func FirebaseMiddleware(firebaseService *FirebaseService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.ErrUnauthorized
			}

			token := authHeader[len("Bearer "):]
			uid, err := firebaseService.VerifyIDToken(c.Request().Context(), token)
			if err != nil {
				return echo.ErrUnauthorized
			}

			c.Set("uid", uid)

			return next(c)
		}
	}
}

func TestMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		fmt.Println("test middleware")

		return next(c)
	}
}

func (f *FirebaseService) VerifyIDToken(ctx context.Context, idToken string) (string, error) {
	token, err := f.app.VerifyIDToken(ctx, idToken)
	if err != nil {
		return "verify token error", err
	}

	uid := token.UID

	return uid, nil
}

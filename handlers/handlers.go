package handlers

import (
	"io"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/estaesta/hijalearn/auth"
	"github.com/estaesta/hijalearn/models"

	"github.com/labstack/echo/v4"
)

func GetProgressUser(c echo.Context, dbClient *firestore.Client) error {
	uid := c.Get("uid").(string)

	iter := dbClient.Collection("users").Doc(uid).Collection("bab")
	iterSnap, err := iter.Documents(c.Request().Context()).GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	dataMap := map[string]interface{}{
		"bab": map[string]interface{}{},
	}
	for _, doc := range iterSnap {
		dataMap["bab"].(map[string]interface{})[doc.Ref.ID] = doc.Data()
	}

	return c.JSON(http.StatusOK, dataMap)
}

func UpdateSubab(c echo.Context, dbClient *firestore.Client) error {
	uid := c.Get("uid").(string)
	bab := c.FormValue("bab")
	subab := c.FormValue("subab")

	progressSubab := map[string]interface{}{
		"subab": map[string]interface{}{
			subab: true,
		},
	}

	doc := dbClient.Collection("users").Doc(uid).Collection("bab").Doc(bab)
	_, err := doc.Set(c.Request().Context(), progressSubab, firestore.MergeAll)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.String(http.StatusOK, uid)
}

func UpdateBab(c echo.Context, dbClient *firestore.Client) error {
	uid := c.Get("uid").(string)
	bab := c.FormValue("bab")

	progressBab := map[string]interface{}{
		"selesai": true,
	}

	doc := dbClient.Collection("users").Doc(uid).Collection("bab").Doc(bab)
	_, err := doc.Set(c.Request().Context(), progressBab, firestore.MergeAll)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.String(http.StatusOK, uid)
}

func UpdateProgressUser(c echo.Context, dbClient *firestore.Client) error {
	if c.FormValue("subab") == "" {
		return UpdateBab(c, dbClient)
	}
	return UpdateSubab(c, dbClient)
}

func InitProgressUser(c echo.Context, dbClient *firestore.Client) error {
	uid := c.Get("uid").(string)
	username := c.FormValue("username")

	newProgress := models.ProgressUser{
		Id:       uid,
		Username: username,
	}

	doc := dbClient.Doc("users/" + uid)
	_, err := doc.Create(c.Request().Context(), newProgress)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	newProgress.Id = ""
	return c.JSON(http.StatusOK, newProgress)
}

func Predict(c echo.Context, url string) error {
	audioFile, err := c.FormFile("audio")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	src, err := audioFile.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	defer src.Close()

	label := c.FormValue("label")

	// send to flask server
	resp, err := http.Post(url, audioFile.Header.Get("Content-Type"), src)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	defer resp.Body.Close()

	// get response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	result := string(body)

	if result == label {
		return c.JSON(http.StatusOK, "benar")
	}
	return c.JSON(http.StatusOK, "salah")
}

func Register(c echo.Context, firebaseService *auth.FirebaseService) error {
	email := c.FormValue("email")
	password := c.FormValue("password")
	username := c.FormValue("username")

	// check if email is already exist
	_, err := firebaseService.GetUserByEmail(c.Request().Context(), email)
	if err == nil {
		return c.JSON(http.StatusBadRequest, "Email is already exist")
	}

	// create user in firebase auth
	user, err := firebaseService.CreateUser(c.Request().Context(), email, password, username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	c.Logger().Info(user)

	// init progress user
	// _, err = dbClient.Collection("users").Doc(user.UID).Create(c.Request().Context(), models.ProgressUser{
	// 	Id:       user.UID,
	// 	Username: username,
	// })

	// auto login
	// token, err := firebaseService.CreateCustomToken(c.Request().Context(), user.UID)
	return c.JSON(http.StatusOK, "User created successfully")
}

func UpdateProfile(c echo.Context, firebaseService *auth.FirebaseService) error {
	// uid := c.Get("uid").(string)
	// username := c.FormValue("username")
	// profilePicture := c.FormFile("profile_picture")

	// TODO; check if profile picture is already exist
	// if exist, delete the old one

	// TODO: upload profile picture to bucket
	// set profile picture url to firebase

	return c.JSON(http.StatusOK, "Profile updated successfully")
}

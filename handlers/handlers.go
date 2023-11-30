package handlers

import (
	"io"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/estaesta/hijalearn/models"

	"github.com/labstack/echo/v4"
)

func GetProgressUser(c echo.Context, dbClient *firestore.Client) error {
	uid := c.Get("uid").(string)

	// db := db.CreateClient(c.Request().Context())
	// defer db.Close()

	// doc := dbClient.Doc("users/" + uid)
	//
	// docSnap, err := doc.Get(c.Request().Context())
	// if err != nil {
	// 	return c.JSON(http.StatusInternalServerError, err)
	// }
	// dataMap := docSnap.Data()
	// fmt.Println(dataMap)

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

	// db := db.CreateClient(c.Request().Context())
	// defer db.Close()

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

	// db := db.CreateClient(c.Request().Context())
	// defer db.Close()

	// babInt, err := strconv.Atoi(bab)
	// if err != nil {
	// 	return c.JSON(http.StatusInternalServerError, err)
	// }
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

	// db := db.CreateClient(c.Request().Context())
	// defer db.Close()

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

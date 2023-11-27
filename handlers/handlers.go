package handlers

import (
	"fmt"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/estaesta/hijalearn/models"

	"github.com/labstack/echo/v4"
)

func GetProgressUser(c echo.Context, dbClient *firestore.Client) error {
	uid := c.Get("uid").(string)

	// db := db.CreateClient(c.Request().Context())
	// defer db.Close()

	doc := dbClient.Doc("progress/" + uid)

	docSnap, err := doc.Get(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	dataMap := docSnap.Data()
	fmt.Println(dataMap)

	return c.JSON(http.StatusOK, dataMap)
}

func UpdateSubab(c echo.Context, dbClient *firestore.Client) error {
	uid := c.Param("uid")
        bab := c.FormValue("bab")
        subab := c.FormValue("subab")

        // db := db.CreateClient(c.Request().Context())
        // defer db.Close()

	progressSubab := models.ProgressSubab{
		Selesai: true,
	}

	doc := dbClient.Collection("users").Doc(uid).Collection("bab").Doc(bab).Collection("subab").Doc(subab)
	wr, err := doc.Set(c.Request().Context(), progressSubab)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	fmt.Println(wr)
        
	return c.String(http.StatusOK, uid)
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

	// doc := db.Doc("progress/" + uid)
	doc := dbClient.Doc("users/" + uid)
	wr, err := doc.Create(c.Request().Context(), newProgress)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	fmt.Println(wr)

	newProgress.Id = ""
	return c.JSON(http.StatusOK, newProgress)
}

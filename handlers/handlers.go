package handlers

import (
	"io"
	"net/http"
	"strconv"

	"cloud.google.com/go/firestore"
	"github.com/estaesta/hijalearn/auth"

	"github.com/labstack/echo/v4"
)

func GetProgressUser(c echo.Context, dbClient *firestore.Client) error {
	uid := c.Get("uid").(string)

	// v1
	// iter := dbClient.Collection("users").Doc(uid).Collection("bab")
	// iterSnap, err := iter.Documents(c.Request().Context()).GetAll()
	// if err != nil {
	// 	return c.JSON(http.StatusInternalServerError, err)
	// }
	//
	// dataMap := map[string]interface{}{
	// 	"bab": map[string]interface{}{},
	// }
	// for _, doc := range iterSnap {
	// 	dataMap["bab"].(map[string]interface{})[doc.Ref.ID] = doc.Data()
	// }

	// v2
	doc := dbClient.Collection("users").Doc(uid)
	docSnap, err := doc.Get(c.Request().Context())
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	dataMap := docSnap.Data()

	dataArray := []interface{}{}
	// add bab id int16
	for babId, babData := range dataMap["module"].(map[string]interface{}) {
		babData.(map[string]interface{})["module_id"] = babId
		dataArray = append(dataArray, babData)
	}

	dataMap["module"] = dataArray

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

func UpdateSubabV2(c echo.Context, dbClient *firestore.Client) error {
	uid := c.Get("uid").(string)
	bab := c.FormValue("bab")
	subab := c.FormValue("subab")

	// v2 change the structure into
	// bab: {
	// "1": {
	// 	"subModuleDone": 1,
	//  "totalSubModule": 3
	// }

	// convert subab to int
	subabInt, err := strconv.Atoi(subab)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	// get progress this subab
	doc := dbClient.Collection("users").Doc(uid).Collection("bab").Doc(bab)
	docSnap, err := doc.Get(c.Request().Context())
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	// check value of subModuleDone
	// if subModuleDone != subab - 1, return error
	subModuleDone := docSnap.Data()["subModuleDone"].(int16)
	if subModuleDone != int16(subabInt-1) {
		return c.JSON(http.StatusBadRequest, "subab is not in order")
	}

	progressSubab := map[string]interface{}{
		"subModuleDone": subabInt,
	}

	_, err = doc.Set(c.Request().Context(), progressSubab, firestore.MergeAll)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	// return this submodule json
	subModuleMap := map[string]interface{}{
		"subModuleDone":  subabInt,
		"totalSubModule": docSnap.Data()["totalSubModule"].(int16),
	}

	return c.JSON(http.StatusOK, subModuleMap)
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
	if c.FormValue("bab") == "" && c.FormValue("subab") == "" {
		return c.JSON(http.StatusBadRequest, "bab or subab is required")
	}
	if c.FormValue("subab") == "" {
		return UpdateBab(c, dbClient)
	}
	return UpdateSubab(c, dbClient)
}

func newModule(totalSubModuleInt int) map[string]interface{} {
	return map[string]interface{}{
		"totalSubModule": totalSubModuleInt,
		"subModuleDone":  0,
		"completed":      false,
	}
}

func InitProgressUser(c echo.Context, dbClient *firestore.Client) error {
	uid := c.Get("uid").(string)

	// doc := dbClient.Collection("users").Doc(uid).Collection("bab").Doc("1")
	// batch.Set(doc, map[string]interface{}{
	// 	"totalSubModule": 28,
	// 	"subModuleDone":  0,
	// }, firestore.MergeAll)
	//
	// doc = dbClient.Collection("users").Doc(uid).Collection("bab").Doc("2")
	// batch.Set(doc, map[string]interface{}{
	// 	"totalSubModule": 28,
	// 	"subModuleDone":  0,
	// }, firestore.MergeAll)
	//
	// _, err = batch.Commit(c.Request().Context())
	// if err != nil {
	// 	c.Logger().Error(err)
	// 	return c.JSON(http.StatusInternalServerError, err)
	// }

	newProgressUser := map[string]interface{}{
		"last_module": 0,
		"module": map[string]interface{}{
			"1": newModule(28),
			"2": newModule(28),
			"3": newModule(28),
		},
	}

	// newProgressUser := map[string]interface{}{
	// 	"last_module": 0,
	// 	"module": []interface{}{
	// 		map[string]interface{}{
	// 			"totalSubModule": totalSubModuleInt,
	// 			"subModuleDone":  0,
	// 			"id":             1,
	// 		},
	// 		map[string]interface{}{
	// 			"totalSubModule": totalSubModuleInt,
	// 			"subModuleDone":  0,
	// 			"id":             2,
	// 		},
	// 	},
	// }

	doc := dbClient.Collection("users").Doc(uid)
	_, err := doc.Set(c.Request().Context(), newProgressUser, firestore.MergeAll)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	return nil
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

func Register(c echo.Context, firebaseService *auth.FirebaseService, dbClient *firestore.Client) error {
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
	// InitProgressUser(c, firebaseService.DbClient)
	err = InitProgressUser(c, dbClient)

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

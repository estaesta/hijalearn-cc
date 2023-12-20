package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"cloud.google.com/go/firestore"
	"github.com/estaesta/hijalearn/auth"

	"github.com/labstack/echo/v4"
)

func GetProgressUser(c echo.Context, dbClient *firestore.Client) error {
	uid := c.Get("uid").(string)

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
		babIdInt, err := strconv.Atoi(babId)
		if err != nil {
			c.Logger().Error(err)
			return c.JSON(http.StatusInternalServerError, err)
		}

		babData.(map[string]interface{})["module_id"] = babIdInt
		dataArray = append(dataArray, babData)
	}
	// sort by module_id
	sort.Slice(dataArray, func(i, j int) bool {
		return dataArray[i].(map[string]interface{})["module_id"].(int) < dataArray[j].(map[string]interface{})["module_id"].(int)
	})

	dataMap["module"] = dataArray

	return c.JSON(http.StatusOK, dataMap)
}

func GetProgressUserModule(c echo.Context, dbClient *firestore.Client) error {
	uid := c.Get("uid").(string)
	moduleId := c.Param("moduleId")

	// v2
	doc := dbClient.Collection("users").Doc(uid)
	docSnap, err := doc.Get(c.Request().Context())
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	dataModule := docSnap.Data()["module"].(map[string]interface{})
	currentModule := dataModule[moduleId].(map[string]interface{})

	// return json of this module data 
	return c.JSON(http.StatusOK, currentModule)
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

	newProgressUser := map[string]interface{}{
		"last_module": 1,
		"module": map[string]interface{}{
			"1": newModule(30),
			"2": newModule(28),
			"3": newModule(28),
			"4": newModule(28),
		},
	}

	doc := dbClient.Collection("users").Doc(uid)
	_, err := doc.Set(c.Request().Context(), newProgressUser, firestore.MergeAll)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	return nil
}

func createTempFile(file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	dst, err := os.CreateTemp("", fmt.Sprintf("%s-*%s", file.Filename, ".wav"))
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}

	return dst.Name(), nil
}

func sendRequest(c echo.Context,filename string, url string, model string) (string, error) {
	form := new(bytes.Buffer)
	writer := multipart.NewWriter(form)
	fw, err := writer.CreateFormFile("file", filepath.Base(filename))
	if err != nil {
		c.Logger().Error(err)
		return "", err
	}
	fd, err := os.Open(filename)
	if err != nil {
		c.Logger().Error(err)
		return "", err
	}
	defer fd.Close()
	_, err = io.Copy(fw, fd)
	if err != nil {
		c.Logger().Error(err)
		return "", err
	}

	formField, err := writer.CreateFormField("model")
	if err != nil {
		c.Logger().Error(err)
		return "", err
	}
	formField.Write([]byte(model))

	writer.Close()

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, form)
	if err != nil {
		c.Logger().Error(err)
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, err := client.Do(req)
	if err != nil {
		c.Logger().Error(err)
		return "", err
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		c.Logger().Error(err)
		return "", err
	}
	// fmt.Printf("%s\n", bodyText)

	return string(bodyText), nil
}

func Predict(c echo.Context, dbClient *firestore.Client, url string) error {
	
	if c.FormValue("caraEja") == "" || c.FormValue("done") == "" {
		return c.JSON(http.StatusBadRequest, "caraEja or done is required")
	}

	// get audio file
	audioFile, err := c.FormFile("audio")
	if err != nil {
		if err == http.ErrMissingFile {
			return c.JSON(http.StatusBadRequest, "audio is required")
		}
		return c.JSON(http.StatusInternalServerError, err)
	}

	label := c.FormValue("caraEja")
	done := c.FormValue("done")
	moduleId := c.FormValue("moduleId")
	if moduleId == "" {
		moduleId = "0"
	}

	// create temp file
	filename, err := createTempFile(audioFile)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	defer os.Remove(filename)

	// send request
	result, err := sendRequest(c, filename, url, moduleId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	// fmt.Println(result)
	// fmt.Println(label)
	// result format
	// {
	// "prediction": "A",
	// "probability": "99.99998807907104"
	// }

	// if result != label {
	// 	return c.JSON(http.StatusOK, "Wrong answer")
	// }

	type Result struct {
		Prediction  string `json:"prediction"`
		Probability string `json:"probability"`
	}

	var resultStruct Result
	err = json.Unmarshal([]byte(result), &resultStruct)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	probability, err := strconv.ParseFloat(resultStruct.Probability, 64)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	fmt.Println(resultStruct.Prediction)
	fmt.Println(label)

	// case insensitive
	if !strings.EqualFold(resultStruct.Prediction, label) || probability < 0.6 {
		response := map[string]interface{}{
			"correct": false,
			"message": "Wrong answer",
			"probability": 0,
		}
		return c.JSON(http.StatusOK, response)
	}

	// if done == "true" {
	// 	// do not update progress
	// 	response := map[string]interface{}{
	// 		"correct": true,
	// 		"probability": probability,
	// 		"message": "Correct answer",
	// 	}
	// 	return c.JSON(http.StatusOK, response)
	// }

	// check progress user
	uid := c.Get("uid").(string)
	fmt.Println(uid)
	doc := dbClient.Collection("users").Doc(uid)
	docSnap, err := doc.Get(c.Request().Context())
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	// dataMap := docSnap.Data()
	// dataModule := dataMap["module"].(map[string]interface{})
	// fmt.Println(moduleId)
	// currentModule := dataModule[moduleId].(map[string]interface{})
	// fmt.Println(currentModule)

	// dataModule = last module
	// lastModule := docSnap.Data()["last_module"].(int64)
	// lastModuleStr := strconv.Itoa(int(lastModule))
	dataModule := docSnap.Data()["module"].(map[string]interface{})
	currentModule := dataModule[moduleId].(map[string]interface{})

	// moduleId := lastModuleStr


	totalSubModule := currentModule["totalSubModule"].(int64)
	subModuleDone := currentModule["subModuleDone"].(int64)
	// check if this is last subModule
	subModuleCompleted := totalSubModule == subModuleDone+1
	fmt.Println(subModuleCompleted)
	fmt.Println(totalSubModule)

	var progressModule map[string]interface{}

	// do not update progress if this module is completed
	if currentModule["completed"].(bool) || done == "true" {
		progressModule = map[string]interface{}{}
	} else {
		progressModule = map[string]interface{}{
			"subModuleDone": subModuleDone + 1,
			"completed":     subModuleCompleted,
		}
	}

	// init progressUser
	var progressUser map[string]interface{}
	// moduleIdInt, err := strconv.Atoi(moduleId)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	moduleIdInt, err := strconv.Atoi(moduleId)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	// update the lastModule to this module
	progressUser = map[string]interface{}{
		"last_module": moduleIdInt,
		"module": map[string]interface{}{
			moduleId: progressModule,
		},
	}

	fmt.Println(progressUser)

	_, err = doc.Set(c.Request().Context(), progressUser, firestore.MergeAll)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	response := map[string]interface{}{
		"correct": true,
		"probability": probability,
		"message": "Correct answer",
	}
	return c.JSON(http.StatusOK, response)
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
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	c.Logger().Info(user)

	// get uid
	uid := user.UID

	// set uid to context
	c.Set("uid", uid)

	err = InitProgressUser(c, dbClient)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

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

package handlers

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"sort"
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
		"last_module": 1,
		"module": map[string]interface{}{
			"1": newModule(30),
			"2": newModule(28),
			"3": newModule(28),
			"4": newModule(28),
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

func sendRequest(c echo.Context,filename string, url string) (string, error) {
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
	// moduleId := c.FormValue("moduleId")

	// create temp file
	filename, err := createTempFile(audioFile)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	defer os.Remove(filename)

	// send request
	result, err := sendRequest(c, filename, url)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	fmt.Println(result)
	fmt.Println(label)

	if result != label {
		return c.JSON(http.StatusOK, "Wrong answer")
	}

	if done == "true" {
		// do not update progress
		return c.JSON(http.StatusOK, "Correct answer")
	}

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
	lastModule := docSnap.Data()["last_module"].(int64)
	lastModuleStr := strconv.Itoa(int(lastModule))
	dataModule := docSnap.Data()["module"].(map[string]interface{})
	currentModule := dataModule[lastModuleStr].(map[string]interface{})

	moduleId := lastModuleStr


	totalSubModule := currentModule["totalSubModule"].(int64)
	subModuleDone := currentModule["subModuleDone"].(int64)
	// check if this module is actually completed
	if currentModule["completed"].(bool) {
		return c.JSON(http.StatusOK, "Correct answer")
	}
	// check if this is last subModule
	subModuleCompleted := totalSubModule == subModuleDone+1
	fmt.Println(subModuleCompleted)
	fmt.Println(totalSubModule)

	// update progress user
	progressModule := map[string]interface{}{
		"subModuleDone": subModuleDone + 1,
		"completed":     subModuleCompleted,
	}

	// init progressUser
	var progressUser map[string]interface{}
	// moduleIdInt, err := strconv.Atoi(moduleId)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	// if this is last subModule, update last_module
	if subModuleCompleted {
		progressUser = map[string]interface{}{
			// "last_module": moduleIdInt + 1,
			"last_module": lastModule + 1,
			"module": map[string]interface{}{
				moduleId: progressModule,
			},
		}
	} else {
		progressUser = map[string]interface{}{
			"module": map[string]interface{}{
				moduleId: progressModule,
			},
		}
	}

	fmt.Println(progressUser)

	_, err = doc.Set(c.Request().Context(), progressUser, firestore.MergeAll)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, "Correct answer")
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

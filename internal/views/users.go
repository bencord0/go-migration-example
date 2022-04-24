package views

import (
	"encoding/json"
	"net/http"

	"github.com/bencord0/go-migration-example/internal/models"
	. "github.com/bencord0/webframework"
)

func GetUsers(db *models.Database) func(*http.Request) *http.Response {
	users := models.UserModel(db)

	return func(req *http.Request) (*http.Response) {
		allUsers, err := users.List()
		if err != nil {
			return JSONResponse(
				map[string]string{"error": err.Error()},
				http.StatusInternalServerError,
				nil,
			)
		}

		return JSONResponse(allUsers, http.StatusOK, nil)
	}
}

func AddUser(db *models.Database) func(*http.Request) *http.Response {
	users := models.UserModel(db)

	return func(req *http.Request) (*http.Response) {
		decoder := json.NewDecoder(req.Body)
		decoder.DisallowUnknownFields()

		var user models.User
		err := decoder.Decode(&user)

		if err != nil {
			return JSONResponse(
				map[string]string{"error": err.Error()},
				http.StatusBadRequest,
				nil,
			)
		}

		err = users.Add(user.Name)
		if err != nil {
			return JSONResponse(
				map[string]string{"error": err.Error()},
				http.StatusInternalServerError,
				nil,
			)
		}

		return JSONResponse(user, http.StatusOK, nil)
	}
}

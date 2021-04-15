package handlers

import (
	"encoding/json"
	"net/http"
	"webmine/models"
	u "webmine/utils"
)

func respond(w http.ResponseWriter, i interface{}, user *models.User) {
	resp := make(map[string]interface{})
	resp["model"] = i
	resp["user"] = user
	u.Respond(w, resp)
}

func getCurrentUser(r *http.Request) *models.User {
	userID := getCurrentUserID(r)
	if userID == 0 {
		return &models.User{}
	}
	user := models.GetUser(userID)
	user.Password = ""
	return user
}

func getCurrentUserID(r *http.Request) uint {
	rawUserID := r.Context().Value("user")
	userID, ok := rawUserID.(uint)
	if !ok {
		return 0
	}
	return userID
}

func convertMap(form map[string][]string) map[string]string {
	mapModel := make(map[string]string)
	for key, value := range form {
		mapModel[key] = u.MakeLine(value)
	}
	return mapModel
}

func decodeForm(form map[string][]string, model interface{}) error {
	mapModel := convertMap(form)
	jsonModel, err := json.Marshal(mapModel)
	if err == nil {
		err = json.Unmarshal(jsonModel, model)
	}
	return err
}

func decodeRequest(r *http.Request, model interface{}) error {
	if r.Header.Get("Content-Type") == "application/json" {
		return json.NewDecoder(r.Body).Decode(model)
	} else {
		r.ParseMultipartForm(2 << 20)
		decodeError := decodeForm(r.Form, model)
		return decodeError
	}
}

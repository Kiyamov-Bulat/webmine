package handlers

import (
	"net/http"
	"webmine/models"
)

func Authenticate(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	if err := decodeRequest(r, user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := user.Login(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	respond(w, nil, user)
}

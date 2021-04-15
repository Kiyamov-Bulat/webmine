package handlers

import (
	"net/http"
	"os"
	m "webmine/models"
	u "webmine/utils"
)

func GetServerInfo(w http.ResponseWriter, r *http.Request) {
	server := m.GetServer()
	if server.ID == "" {
		if err := server.UpdateId(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	if err := server.UpdateState(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := server.UpdateUptime(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respond(w, server, getCurrentUser(r))
}

func ReloadServer(w http.ResponseWriter, r *http.Request) {
	server := m.GetServer()
	if server.ID == "" {
		if err := server.UpdateId(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	if err := server.Reload(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respond(w, server, getCurrentUser(r))
}

func GetServerDataArchive(w http.ResponseWriter, r *http.Request) {
	err := u.RecursiveZip(m.SERVER_DATA_DIR_PATH, m.SERVER_DATA_ZIP_PATH)
	if err != nil {
		http.Error(w, "failed to zip archive", http.StatusInternalServerError)
		return
	}
	u.ServeFile(w, r, m.SERVER_DATA_ZIP_PATH)
	defer os.Remove(m.SERVER_DATA_ZIP_PATH)
}

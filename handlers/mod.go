package handlers

import (
	"net/http"
	"path"
	"strconv"
	m "webmine/models"
	u "webmine/utils"

	"github.com/gorilla/mux"
)

func GetAllModNames(w http.ResponseWriter, r *http.Request) {
	modNames := m.GetAllModNames()
	user := getCurrentUser(r)
	respond(w, modNames, user)
}

func UploadMod(w http.ResponseWriter, r *http.Request) {
	mod := &m.Mod{}
	if err := decodeRequest(r, mod); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	mod.SaveFile()
	mod.Content = nil
	respond(w, mod, getCurrentUser(r))
}

func DownLoadMod(w http.ResponseWriter, r *http.Request) {
	var mod *m.Mod

	smodID := mux.Vars(r)["modID"]
	modID, err := u.Stou(smodID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	mod = m.GetMod(modID)
	if !mod.IsValid() {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	filePath := path.Join(m.MOD_DIR_PATH, mod.Name)
	w.Header().Set("Content-Disposition", "attachment; filename="+strconv.Quote(mod.Name))
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeFile(w, r, filePath)
}

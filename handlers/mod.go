package handlers

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"time"
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
	r.ParseMultipartForm(32 << 20)
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	mod.Content = buf.Bytes()
	mod.Name = header.Filename
	mod.ModTime = time.Now()
	mod.SaveFile()
	mod.Create()
	mod.Content = nil
	respond(w, mod, getCurrentUser(r))
}

func DownLoadMod(w http.ResponseWriter, r *http.Request) {
	mod := getModFromRequest(r)
	if !mod.IsValid() {
		http.Error(w, "mod's id is not valid", http.StatusBadRequest)
		return
	}
	filePath := mod.GetFullFilePath()
	u.ServeFile(w, r, filePath)
}

func DeleteMod(w http.ResponseWriter, r *http.Request) {
	mod := getModFromRequest(r)
	if !mod.IsValid() {
		http.Error(w, "mod's id is not valid", http.StatusBadRequest)
		return
	}
	mod.Delete()
	mod.Content = nil
	respond(w, mod, getCurrentUser(r))
}

func GetModsArchive(w http.ResponseWriter, r *http.Request) {
	err := u.RecursiveZip(m.MOD_DIR_PATH, m.MOD_ZIP_PATH)
	if err != nil {
		http.Error(w, "failed to zip archive", http.StatusInternalServerError)
		return
	}
	u.ServeFile(w, r, m.MOD_ZIP_PATH)
	defer os.Remove(m.MOD_ZIP_PATH)
}

func getModFromRequest(r *http.Request) *m.Mod {
	smodID := mux.Vars(r)["modID"]
	modID, err := u.Stou(smodID)
	if err != nil {
		return nil
	}
	return m.GetMod(modID)
}

package main

import (
	"log"
	"net/http"
	"os"

	mwr "webmine/auth"
	h "webmine/handlers"
	u "webmine/utils"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.PathPrefix("/frontend/").Handler(http.StripPrefix("/frontend/dist/", http.FileServer(http.Dir("./frontend/dist"))))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	r.HandleFunc("/api/login", h.Authenticate).Methods("POST")
	r.HandleFunc("/api/server", h.GetServerInfo).Methods("GET")
	r.HandleFunc("/api/server/reload", h.ReloadServer).Methods("GET")
	r.HandleFunc("/api/mods/{modID:[0-9]+}", h.DownLoadMod).Methods("GET")
	r.HandleFunc("/api/mods/{modID:[0-9]+}", h.DeleteMod).Methods("POST")
	r.HandleFunc("/api/archive/mods", h.GetModsArchive).Methods("GET")
	r.HandleFunc("/api/archive/data", h.GetServerDataArchive).Methods("GET")
	r.HandleFunc("/api/mods", h.GetAllModNames).Methods("GET")
	r.HandleFunc("/api/mods", h.UploadMod).Methods("POST")
	r.HandleFunc("/", h.IndexHandler).Methods("GET")
	r.Use(mwr.JwtAuth)
	port := ":" + os.Getenv("PORT")
	if port == ":" {
		port += u.DefaultEnv["PORT"]
	}
	log.Println(port)
	log.Fatal(http.ListenAndServe(port, r))
}

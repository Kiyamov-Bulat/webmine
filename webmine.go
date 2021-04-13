package main

import (
	"log"
	"net/http"
	"os"

	mwr "webmine/auth"
	h "webmine/handlers"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.PathPrefix("/frontend/").Handler(http.StripPrefix("/frontend/dist/", http.FileServer(http.Dir("./frontend/dist"))))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	r.HandleFunc("/api/login", h.Authenticate).Methods("POST")
	r.HandleFunc("/api/mods/{modID:[0-9]+}", h.DownLoadMod).Methods("GET")
	r.HandleFunc("/api/mods", h.GetAllModNames).Methods("GET")
	r.HandleFunc("/", h.IndexHandler).Methods("GET")
	r.Use(mwr.JwtAuth)
	port := ":" + os.Getenv("PORT")
	log.Println(port)
	log.Fatal(http.ListenAndServe(port, r))
}

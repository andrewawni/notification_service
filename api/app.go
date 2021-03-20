package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type App struct {
	router *mux.Router
}

func (app *App) Init() {
	app.router = mux.NewRouter()
	app.drawRoutes()
}

func (app *App) Run(address string) {
	log.Fatal(http.ListenAndServe(address, app.router))
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func (app *App) drawRoutes() {
	app.router.HandleFunc("/sendNotification", app.sendNotificationHandler).Methods("POST")
}

func (app *App) sendNotificationHandler(w http.ResponseWriter, r *http.Request) {
}

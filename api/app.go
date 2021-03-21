package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/andrewawni/notification_service/common/messaging"
	"github.com/andrewawni/notification_service/common/models"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type App struct {
	router          *mux.Router
	messagingClient messaging.IMessagingClient
}

const singleNotificationQueueName = "notification_service:assemble_single_notifications_jobs"
const groupNotificationQueueName = "notification_service:assemble_group_notifications_jobs"

func (app *App) Init() {
	app.router = mux.NewRouter()
	app.drawRoutes()
	app.messagingClient = &messaging.MessagingClient{}
	app.messagingClient.ConnectToBroker(os.Getenv("RABBITMQ_URL"))
	log.Printf("connected to rabbitmq")
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
	app.router.HandleFunc("/", app.index).Methods("POST")
	app.router.HandleFunc("/sendSingleNotification", app.sendSingleNotificationHandler).Methods("POST")
	app.router.HandleFunc("/sendGroupNotification", app.sendGroupNotificationHandler).Methods("POST")
}

func (app *App) index(w http.ResponseWriter, r *http.Request) {
	log.Printf("I'm working")
}

func (app *App) sendSingleNotificationHandler(w http.ResponseWriter, r *http.Request) {
	id := uuid.New().String()
	body := models.SingleNotification{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid payload")
		return
	}
	body.ID = id
	payload, _ := json.Marshal(&body)
	app.messagingClient.PublishOnQueue([]byte(payload), singleNotificationQueueName)
	respondWithJSON(w, http.StatusAccepted, body)
}

func (app *App) sendGroupNotificationHandler(w http.ResponseWriter, r *http.Request) {
	id := uuid.New().String()
	body := models.GroupNotification{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid payload")
		return
	}
	body.ID = id
	payload, _ := json.Marshal(&body)
	app.messagingClient.PublishOnQueue([]byte(payload), groupNotificationQueueName)
	respondWithJSON(w, http.StatusAccepted, body)
}

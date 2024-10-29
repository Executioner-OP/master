package httpServer

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Executioner-OP/master/pkg/dbHandler"
	"github.com/Executioner-OP/master/pkg/queueHandler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type RequestData struct {
	ID   string `json:"id"`
	Code string `json:"code"`
}

func Init() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Define the main endpoint to handle incoming requests
	r.Post("/request", handleRequest)

	fmt.Println("Connected to HTTP server on :3000")
	http.ListenAndServe(":3000", r)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	var requestData RequestData

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Add request to the database
	id, err := dbHandler.AddToDb(requestData.Code, false)

	requestData.ID = id

	if err != nil {
		http.Error(w, "Failed to add request to database", http.StatusInternalServerError)
		return
	}

	queueHandler.AddToQueue(id)

	// Send response
	w.Header().Set("Content-Type", "application/json")

}

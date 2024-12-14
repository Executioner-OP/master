package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Executioner-OP/master/db"
	"github.com/Executioner-OP/master/queue"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Init() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Define the main endpoint to handle incoming requests
	r.Post("/request", handleRequest)

	fmt.Println("Connected to HTTP server on :3000")
	http.ListenAndServe(":3000", r)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	var requestData db.ExecutionRequest

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Add request to the database
	requestData, err := db.AddToDb(requestData)

	if err != nil {
		http.Error(w, "Failed to add request to database", http.StatusInternalServerError)
		return
	}
	fmt.Println("Added request to database")

	// encoding the request data to bytes and adding it to the queue
	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(requestData)
	err = queue.AddToQueue(reqBodyBytes.Bytes())

	if err != nil {
		http.Error(w, "Failed to add request to queue", http.StatusInternalServerError)
		return
	}
	fmt.Println("Added request to queue")

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "Request added to queue"}`))
}

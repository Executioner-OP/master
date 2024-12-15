package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Executioner-OP/master/db"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func InitHttpServer(taskChannel chan db.ExecutionRequest) {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Binding Channel
	TaskChannel = taskChannel

	// Define the main endpoint to handle incoming requests
	r.Post("/request", handleExecutionRequest)
	r.Post("/getSubmission", getSubmission)

	fmt.Println("Connected to HTTP server on :3000")
	http.ListenAndServe(":3000", r)
}

func handleExecutionRequest(w http.ResponseWriter, r *http.Request) {
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

	// Add request to channel
	TaskChannel <- requestData

	// encoding the request data to bytes and adding it to the queue
	// reqBodyBytes := new(bytes.Buffer)
	// json.NewEncoder(reqBodyBytes).Encode(requestData)
	// err = queue.AddToQueue(reqBodyBytes.Bytes())

	// if err != nil {
	// 	http.Error(w, "Failed to add request to queue", http.StatusInternalServerError)
	// 	return
	// }
	fmt.Println("Added request to queue")

	// Send response
	if err := json.NewEncoder(w).Encode(requestData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getSubmission(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Token string `json:"token"`
	}

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Convert token string to ObjectID
	objectID, err := primitive.ObjectIDFromHex(req.Token)
	if err != nil {
		http.Error(w, "Invalid ObjectID format", http.StatusBadRequest)
		return
	}

	// Retrieve the object from the database
	result, err := db.ReadFromDb(objectID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading value from DB: %v", err), http.StatusInternalServerError)
		return
	}
	fmt.Println(result)
	// Encode the result and send it back in the response
	// w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

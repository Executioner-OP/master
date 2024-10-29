package httpServer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/Executioner-OP/master/pkg/dbHandler"
	"github.com/Executioner-OP/master/pkg/queueHandler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

type RequestData struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func Init() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Define the main endpoint to handle incoming requests
	r.Post("/request", handleRequestSync)

	fmt.Println("Connected to HTTP server on :3000")
	http.ListenAndServe(":3000", r)
}

// Synchronous handler to ensure DB and queue are updated before responding
func handleRequestSync(w http.ResponseWriter, r *http.Request) {
	var requestData RequestData

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Generate a new UUID for the request
	requestData.ID = uuid.New().String()

	// Use WaitGroup to manage concurrency
	var wg sync.WaitGroup
	wg.Add(2) // We have two concurrent tasks

	// Channels to capture errors
	dbErrChan := make(chan error, 1)
	queueErrChan := make(chan error, 1)

	// Run DB operation in a separate goroutine
	go func() {
		defer wg.Done()
		if err := dbHandler.AddToDB(requestData); err != nil {
			dbErrChan <- err
		}
	}()

	// Run Queue operation in a separate goroutine
	go func() {
		defer wg.Done()
		if err := queueHandler.AddToQueue(requestData); err != nil {
			queueErrChan <- err
		}
	}()

	// Wait for both goroutines to finish
	wg.Wait()
	close(dbErrChan)
	close(queueErrChan)

	// Check for errors
	if dbErr := <-dbErrChan; dbErr != nil {
		http.Error(w, "Failed to add to database: "+dbErr.Error(), http.StatusInternalServerError)
		return
	}
	if queueErr := <-queueErrChan; queueErr != nil {
		http.Error(w, "Failed to add to queue: "+queueErr.Error(), http.StatusInternalServerError)
		return
	}

	// If both succeeded, respond to client
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Request processed successfully"))
}

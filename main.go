package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"

	"github.com/Executioner-OP/master/db"
	"github.com/Executioner-OP/master/queue"
	"github.com/Executioner-OP/master/server"
)

func main() {

	// Find .evn
	fmt.Println("Program Started")
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
	RABBITMQ_URI := os.Getenv("RABBITMQ_URI")
	MONGO_URI := os.Getenv("MONGO_URI")

	// Initialize database and queue synchronously
	if err := db.Connect(MONGO_URI); err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	if err := queue.Init(RABBITMQ_URI); err != nil {
		log.Fatalf("Failed to initialize queue: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	// Start HTTP server concurrently
	go func() {
		defer wg.Done()
		server.InitHttpServer()
	}()

	go func() {
		defer wg.Done()
		server.InitGrpcServer()
	}()

	// Wait fot both servers to finish
	wg.Wait()
}

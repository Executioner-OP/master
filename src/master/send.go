package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/Executioner-OP/master/pkg/dbHandler"
	"github.com/Executioner-OP/master/pkg/httpServer"
	"github.com/Executioner-OP/master/pkg/queueHandler"
)

func main() {

	// Find .evn
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
	RABBITMQ_URI := os.Getenv("RABBITMQ_URI")
	MONGO_URI := os.Getenv("MONGO_URI")

	// Initialize database and queue synchronously
	if err := dbHandler.Connect(MONGO_URI); err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	if err := queueHandler.Init(RABBITMQ_URI); err != nil {
		log.Fatalf("Failed to initialize queue: %v", err)
	}

	httpServer.Init()
}

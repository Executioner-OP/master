package dbHandler

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	Connect()
}

func Connect() *mongo.Collection {
	// Find .evn
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	// Get value from .env
	MONGO_URI := os.Getenv("MONGO_URI")

	// Connect to the database.
	clientOption := options.Client().ApplyURI(MONGO_URI)
	client, err := mongo.Connect(context.Background(), clientOption)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection.
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	// Create collection
	collection := client.Database("testdb").Collection("test")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to db")

	return collection
}

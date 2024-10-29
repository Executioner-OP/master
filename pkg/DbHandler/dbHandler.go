package dbHandler

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ExecutionRequest struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	Code   string             `bson:"code,omitempty"`
	Output string             `bson:"output,omitempty"`
	IsDone bool               `bson:"isDone"`
}

var (
	collection *mongo.Collection
)

// Helper function to handle errors
func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func Connect(MONGO_URI string) error {
	clientOption := options.Client().ApplyURI(MONGO_URI)
	client, err := mongo.Connect(nil, clientOption)
	failOnError(err, "Failed to connect to MongoDB")

	err = client.Ping(nil, nil)
	failOnError(err, "Failed to ping MongoDB")

	collection = client.Database("executorDB").Collection("executions")
	fmt.Println("Connected to db")

	return nil
}

func AddToDb(code string, IsDone bool) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	request := ExecutionRequest{
		Code:   code,
		IsDone: IsDone,
	}

	// Insert document and get the generated ID
	result, err := collection.InsertOne(ctx, request)
	if err != nil {
		return "", fmt.Errorf("failed to insert execution request: %w", err)
	}

	// Get the inserted ID and convert to string
	id := result.InsertedID.(primitive.ObjectID).Hex()

	return id, nil
}

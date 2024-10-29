package dbHandler

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ExecutionRequest struct {
	ID     string `bson:"id"`
	Code   string `bson:"code,omitempty"`
	Output string `bson:"output,omitempty"`
	IsDone bool   `bson:"isDone"`
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

	// Connect to the database.
	clientOption := options.Client().ApplyURI(MONGO_URI)
	client, err := mongo.Connect(context.Background(), clientOption)
	failOnError(err, "Failed to connect to MongoDB")

	// Check the connection.
	err = client.Ping(context.Background(), nil)
	failOnError(err, "Failed to ping MongoDB")

	// Create collection
	collection = client.Database("executorDB").Collection("executions")
	failOnError(err, "Failed to create collection")

	fmt.Println("Connected to db")

	return nil
}

func AddExecutionRequest(id string, code string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	request := ExecutionRequest{
		ID:     id,
		Code:   code,
		IsDone: false,
	}

	_, err := collection.InsertOne(ctx, request)
	if err != nil {
		return "", fmt.Errorf("failed to insert execution request: %w", err)
	}

	return id, nil
}

// UpdateExecutionRequest updates the execution request with output and marks it as done
func UpdateExecutionRequest(id string, output string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"id": id}
	update := bson.M{
		"$set": bson.M{
			"output": output,
			"isDone": true,
		},
	}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update execution request: %w", err)
	}

	return nil
}

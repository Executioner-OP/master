package db

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
	ID             primitive.ObjectID `bson:"_id"`
	Code           string             `bson:"code"`
	IsDone         bool               `bson:"isDone"`
	LanguageId     int                `bson:"languageId"`
	StandardInput  string             `bson:"standardInput"`
	StandardOutput string             `bson:"standardOutput"`
	ExpectedOutput string             `bson:"expectedOutput"`
	Status         string             `bson:"status"`
	Verdict        string             `bson:"verdict"`
	TimeLimit      int                `bson:"timeLimit"`
	MemoryLimit    int                `bson:"memoryLimit"`
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

	// err = client.Ping(nil, nil)
	// failOnError(err, "Failed to ping MongoDB")

	collection = client.Database("executorDB").Collection("executions")
	fmt.Println("Connected to db")

	return nil
}

func AddToDb(req ExecutionRequest) (ExecutionRequest, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req.ID = primitive.NewObjectID()

	// Insert document and get the generated ID
	_, err := collection.InsertOne(ctx, req)
	if err != nil {
		return req, fmt.Errorf("failed to insert execution request: %w", err)
	}
	fmt.Println("Execution Request: ", req)
	return req, nil
}

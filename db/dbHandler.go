package db

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
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

	// Encode fields to base64
	var err error
	req.Code, err = encodeBase64(req.Code)
	if err != nil {
		return req, fmt.Errorf("failed to encode Code: %w", err)
	}
	req.StandardInput, err = encodeBase64(req.StandardInput)
	if err != nil {
		return req, fmt.Errorf("failed to encode StandardInput: %w", err)
	}
	req.StandardOutput, err = encodeBase64(req.StandardOutput)
	if err != nil {
		return req, fmt.Errorf("failed to encode StandardOutput: %w", err)
	}
	req.ExpectedOutput, err = encodeBase64(req.ExpectedOutput)
	if err != nil {
		return req, fmt.Errorf("failed to encode ExpectedOutput: %w", err)
	}

	// Insert document and get the generated ID
	_, err = collection.InsertOne(ctx, req)
	if err != nil {
		return req, fmt.Errorf("failed to insert execution request: %w", err)
	}
	fmt.Println("Execution Request: ", req)
	return req, nil
}

func encodeBase64(data string) (string, error) {
	encoded := base64.StdEncoding.EncodeToString([]byte(data))
	return encoded, nil
}

func ReadFromDb(req primitive.ObjectID) (ExecutionRequest, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Get document from Db
	var result ExecutionRequest
	err := collection.FindOne(ctx, bson.M{"_id": req}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ExecutionRequest{}, fmt.Errorf("no document found with the given ID")
		}
		return ExecutionRequest{}, fmt.Errorf("failed to find document: %w", err)
	}

	// Decode base64-encoded fields
	result.Code, err = decodeBase64(result.Code)
	if err != nil {
		return ExecutionRequest{}, fmt.Errorf("failed to decode Code: %w", err)
	}
	result.StandardInput, err = decodeBase64(result.StandardInput)
	if err != nil {
		return ExecutionRequest{}, fmt.Errorf("failed to decode StandardInput: %w", err)
	}
	result.StandardOutput, err = decodeBase64(result.StandardOutput)
	if err != nil {
		return ExecutionRequest{}, fmt.Errorf("failed to decode StandardOutput: %w", err)
	}
	result.ExpectedOutput, err = decodeBase64(result.ExpectedOutput)
	if err != nil {
		return ExecutionRequest{}, fmt.Errorf("failed to decode ExpectedOutput: %w", err)
	}

	log.Println(result)
	return result, nil
}

func decodeBase64(encoded string) (string, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	return string(decodedBytes), nil
}

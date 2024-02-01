package models

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB struct holds the MongoDB client and provides methods to interact with MongoDB.
type MongoDB struct {
	client *mongo.Client
}

// NewMongoDB creates a new MongoDB instance with the provided connection string.
func NewMongoDB() (*MongoDB, error) {
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	mongoClient, err := mongo.Connect(
		ctx,
		options.Client().ApplyURI("mongodb://inwza:strongpassword@localhost:27017/"),
	)

	if err != nil {
		log.Fatalf("connection error :%v", err)
		return nil, err
	}

	return &MongoDB{client: mongoClient}, nil
}

// ExampleMethod demonstrates a simple MongoDB operation (e.g., inserting a document).
func (m *MongoDB) InitFirstDB() error {
	collection := m.client.Database("first_db").Collection("Test_collection")

	document := bson.D{
		{Key: "key1", Value: "value1"},
		{Key: "key2", Value: "value2"},
	}

	insertResult, err := collection.InsertOne(context.Background(), document)
	if err != nil {
		return err
	}

	fmt.Printf("Inserted document with ID %v\n", insertResult.InsertedID)
	return nil
}

func (m *MongoDB) GetFirstDB() ([]bson.M, error) {
	collection := m.client.Database("first_db").Collection("Test_collection")

	// Create a filter for all documents
	filter := bson.D{{}}

	// Find all documents in the collection
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	// Decode the results into a slice of bson.M
	var results []bson.M
	if err := cursor.All(context.Background(), &results); err != nil {
		return nil, err
	}

	return results, nil
}

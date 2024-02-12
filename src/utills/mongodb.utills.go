package utills

import (
	"context"
	"fmt"
	"log"
	"petpal-backend/src/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewMongoDB creates a new MongoDB instance with the provided connection string.
func NewMongoDB() (*models.MongoDB, error) {
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	//"mongodb://inwza:strongpassword@localhost:27017/"
	mongoClient, err := mongo.Connect(
		ctx,
		options.Client().ApplyURI("mongodb://inwza:strongpassword@localhost:27017/"),
	)

	if err != nil {
		log.Fatalf("connection error :%v", err)
		return nil, err
	}

	return &models.MongoDB{Client: mongoClient, DbName: "petpal"}, nil
}

func GetFirstDB(m *models.MongoDB) ([]bson.M, error) {
	// Create a filter for all documents
	filter := bson.D{{}}
	results, _ := m.Read("collection1", filter)

	return results, nil
}

// ExampleMethod demonstrates a simple MongoDB operation (e.g., inserting a document).
func AddMockDataToDB(m *models.MongoDB) error {
	// Example of using co
	collection := m.Collection("collection1")

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

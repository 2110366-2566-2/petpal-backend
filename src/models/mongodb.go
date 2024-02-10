package models

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoDB struct holds the MongoDB client and provides methods to interact with MongoDB.
type MongoDB struct {
	Client *mongo.Client
	DbName string
}

func (m *MongoDB) Collection(collectionName string) *mongo.Collection {
	collection := m.Client.Database(m.DbName).Collection(collectionName)
	return collection
}

func (m *MongoDB) Read(collection_name string, filter bson.D) ([]bson.M, error) {
	collection := m.Collection(collection_name)

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

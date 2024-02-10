package utills

import (
	"context"
	"petpal-backend/src/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetSVCPs(db *models.MongoDB, filter bson.D, page int64, per int64) ([]models.SVCP, error) {
	collection := db.Collection("svcp")

	// define options for pagination
	opts := options.Find().SetSkip(page * per).SetLimit(per)

	// Find all documents in the collection
	cursor, err := collection.Find(context.Background(), filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	// iterate over the cursor and decode each document
	var svcps []models.SVCP
	if err := cursor.All(context.Background(), &svcps); err != nil {
		return nil, err
	}

	return svcps, err
}

func GetSVCPByID(db *models.MongoDB, id string) (*models.SVCP, error) {
	// get collection
	collection := db.Collection("svcp")

	// find service provider by *SVCPID*, could be changed to individual id when it exists
	var svcp models.SVCP = models.SVCP{}
	filter := bson.D{{Key: "SVCPID", Value: id}}
	err := collection.FindOne(context.Background(), filter).Decode(&svcp)
	if err != nil {
		return nil, err
	}

	return &svcp, nil
}
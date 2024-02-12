package utills

import (
	"context"
	"petpal-backend/src/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InsertSVCP(db *models.MongoDB, svcp *models.SVCP) (*models.SVCP, error) {
	// Get the users collection
	collection := db.Collection("svcp")

	// Insert the user into the collection
	_, err := collection.InsertOne(context.Background(), svcp)
	if err != nil {
		return nil, err
	}

	// Return the inserted user
	return svcp, nil
}

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
func GetSVCPByEmail(db *models.MongoDB, email string) (*models.SVCP, error) {
	// get collection
	collection := db.Collection("svcp")

	// find service provider by *SVCPID*, could be changed to individual id when it exists
	var svcp models.SVCP = models.SVCP{}
	filter := bson.D{{Key: "SVCPEmail", Value: email}}
	err := collection.FindOne(context.Background(), filter).Decode(&svcp)
	if err != nil {
		return nil, err
	}

	return &svcp, nil
}

func UpdateSVCP(db *models.MongoDB, id string, svcp models.SVCP) error {
	// get collection
	collection := db.Collection("svcp")

	// update service provider by *SVCPID*, could be changed to individual id when it exists
	filter := bson.D{{Key: "SVCPID", Value: id}}
	update := bson.D{{Key: "$set", Value: svcp}}
	res, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	if res.ModifiedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func ChangePassword(email string, newPassword string, db *models.MongoDB) (string, error) {
	// get collection
	svcp_collection := db.Collection("user")

	// find user by id
	var svcp models.User = models.User{}
	filter := bson.D{{Key: "email", Value: email}}
	err := svcp_collection.FindOne(context.Background(), filter).Decode(&svcp)
	if err != nil {
		return "SVCP not found (email=" + email + ")", err
	}

	// update user with new default bank account
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "password", Value: newPassword},
		}},
	}
	_, err = svcp_collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return "Failed to update service provider password", err
	}

	return "", nil
}

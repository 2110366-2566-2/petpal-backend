package utills

import (
	"context"
	"encoding/base64"
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

func UpdateSVCP(db *models.MongoDB, id string, svcp *bson.M) error {
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

func EditDescription(db *models.MongoDB, email string, description string) error {
	// get collection
	collection := db.Collection("svcp")

	// update service provider by *SVCPID*, could be changed to individual id when it exists
	filter := bson.D{{Key: "SVCPEmail", Value: email}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "description", Value: description}}}}
	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return nil
}
// For upload license file to SVCP Collection in mongoDB
func UploadSVCPLicense(db *models.MongoDB, fileContent []byte, SVCPEmail string) error {
	// Encode the file content to base64 string
	encodedFileContent := base64.StdEncoding.EncodeToString(fileContent)
	// get collection
	svcpCollection := db.Collection("svcp")
	// Update the license field in "scvp" collection
	filter := bson.D{{"SVCPEmail", SVCPEmail}}
	update := bson.D{{"$set", bson.D{{"license", encodedFileContent}}}}
	// Updates the first document that has the specified "SVCPUsername" value
	_, err := svcpCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func AddService(db *models.MongoDB, email string, service models.Service) error {
	svcp_collection := db.Collection("svcp")

	svcp := models.SVCP{}
	filter := bson.D{{Key: "SVCPEmail", Value: email}}
	err := svcp_collection.FindOne(context.Background(), filter).Decode(&svcp)
	if err != nil {
		return err
	}

	// append service to services
	svcp.Services = append(svcp.Services, service)

	// update service provider
	filter = bson.D{{Key: "SVCPEmail", Value: email}}
	update := bson.D{{Key: "$set", Value: svcp}}
	res, err := svcp_collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	if res.ModifiedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func DeleteBankAccount(db *models.MongoDB, email string) error {
	// get collection
	svcp_collection := db.Collection("svcp")

	// find service provider by id
	filter := bson.D{{Key: "SVCPEmail", Value: email}}
	result := svcp_collection.FindOne(context.Background(), filter)
	if result.Err() != nil {
		return result.Err()
	}

	// update default bank account to empty
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "defaultAccountNumber", Value: ""},
			{Key: "defaultBank", Value: ""},
		}},
	}

	_, err := svcp_collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	
	return nil
}

func SetDefaultBankAccount(email string, defaultAccountNumber string, defaultBank string, db *models.MongoDB) (string, error) {
	// get collection
	svcp_collection := db.Collection("svcp")

	// find service provider by id
	filter := bson.D{{Key: "SVCPEmail", Value: email}}
	result := svcp_collection.FindOne(context.Background(), filter)
	if result.Err() != nil {
		return "Service provider not found", result.Err()
	}

	// update service provider with new default bank account
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "defaultAccountNumber", Value: defaultAccountNumber},
			{Key: "defaultBank", Value: defaultBank},
		}},
	}

	_, err := svcp_collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return "Error updating service provider", err
	}
	
	return "", nil
}
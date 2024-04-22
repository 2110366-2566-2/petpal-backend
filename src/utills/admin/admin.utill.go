package admin

import (
	"context"
	"errors"
	"petpal-backend/src/models"

	"go.mongodb.org/mongo-driver/bson"
)

func GetAdminByEmail(db *models.MongoDB, email string) (*models.Admin, error) {
	// get collection
	collection := db.Collection("admin")

	// find admin by email
	var admin models.Admin = models.Admin{}
	filter := bson.D{{Key: "email", Value: email}}
	err := collection.FindOne(context.Background(), filter).Decode(&admin)
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func InsertAdmin(db *models.MongoDB, admin *models.Admin) (*models.Admin, error) {
	// get collection
	collection := db.Collection("admin")

	// insert admin
	_, err := collection.InsertOne(context.Background(), admin)
	if err != nil {
		return nil, err
	}
	return admin, nil
}

func AdminUpdateSVCPVerify(db *models.MongoDB, svcpID string, verify bool) error {
	// get collection
	collection := db.Collection("svcp")

	// update SVCP verified status
	filter := bson.D{{Key: "SVCPID", Value: svcpID}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "isVerified", Value: verify}}}}
	res, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("no matched SVCP found")
	}
	return nil
}

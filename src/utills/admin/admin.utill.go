package admin

import (
	"context"
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
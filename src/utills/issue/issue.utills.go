package utills

import (
	"context"
	"petpal-backend/src/models"
)

func CreateIssue(db *models.MongoDB, issue *models.CreateIssue) error {
	collection := db.Collection("issue")

	_, err := collection.InsertOne(context.Background(), issue)
	if err != nil {
		return err
	}

	return nil
}

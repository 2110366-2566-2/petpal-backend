package utills

import (
	"context"
	"errors"
	"petpal-backend/src/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateIssue(db *models.MongoDB, issue *models.CreateIssue) error {
	collection := db.Collection("issue")

	_, err := collection.InsertOne(context.Background(), issue)
	if err != nil {
		return err
	}

	return nil
}

func GetIssues(db *models.MongoDB, filter bson.M, page int64, per int64) ([]models.Issue, error) {
	collection := db.Collection("issue")

	opts := options.Find().SetSkip(page * per).SetLimit(per)

	cursor, err := collection.Find(context.Background(), filter, opts)
	if err != nil {
		return nil, err
	}

	var issues []models.Issue
	err = cursor.All(context.Background(), &issues)
	if err != nil {
		return nil, err
	}

	return issues, nil
}

func AdminAcceptIssue(db *models.MongoDB, issueID string, adminID string) error {
	collection := db.Collection("issue")

	oid, _ := primitive.ObjectIDFromHex(issueID)
	filter := bson.M{"_id": oid}
	update := bson.M{"$set": bson.M{"workingAdminID": adminID}}

	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("no issue found with the given ID")
	}

	return nil
}

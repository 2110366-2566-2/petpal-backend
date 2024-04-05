package utills

import (
	"context"
	"errors"
	"petpal-backend/src/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateIssue(db *models.MongoDB, issue *models.CreateIssue) error {
	collection := db.Collection("issue")

	new_issue := models.Issue{
		IssueID:             primitive.NewObjectID().Hex(),
		IssueDate:           time.Now(),
		IsResolved:          false,
		WorkingAdminID:      "",
		ReporterID:          issue.ReporterID,
		ReporterType:        issue.ReporterType,
		Details:             issue.Details,
		AttachedImg:         issue.AttachedImg,
		IssueType:           issue.IssueType,
		AssociatedBookingID: issue.AssociatedBookingID,
	}

	_, err := collection.InsertOne(context.Background(), new_issue)
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

	filter := bson.M{"_id": issueID}

	issue_to_update := models.Issue{}
	err := collection.FindOne(context.Background(), filter).Decode(&issue_to_update)
	if err != nil {
		return err
	}

	if issue_to_update.WorkingAdminID != "" {
		return errors.New("issue already accepted by another admin")
	}

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

func AdminResolveIssue(db *models.MongoDB, issueID string, adminID string) error {
	collection := db.Collection("issue")

	filter := bson.M{"_id": issueID, "workingAdminID": adminID}
	update := bson.M{"$set": bson.M{"isResolved": true, "resolveDate": time.Now()}}
	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("no issue found with the given ID and workingAdminID")
	}

	return nil
}

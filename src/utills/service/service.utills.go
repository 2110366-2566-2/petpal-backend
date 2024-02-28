package utills

import (
	"context"
	"petpal-backend/src/models"

	"go.mongodb.org/mongo-driver/bson"
)

func UpdateFeedbackToService(db *models.MongoDB, service_id string, user_id string, feedback models.Feedback) error {
	// add feedback to service
	booking_collection := db.Client.Database("petpal").Collection("booking")
	filter := bson.M{"serviceID": service_id, "userID": user_id}
	update := bson.M{"$set": bson.M{"feedback": feedback}}
	_, err := booking_collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		// show error message
		return err
	}

	collention_temp, err := booking_collection.Aggregate(context.Background(), bson.A{
		bson.M{"$unwind": "$feedback"},
		bson.M{"$match": bson.M{"serviceID": service_id}},
		bson.M{"$group": bson.M{"_id": "$serviceID", "averageRating": bson.M{"$avg": "$feedback.rating"}}},
	})
	if err != nil {
		return err
	}
	var rating_result []bson.M
	if err = collention_temp.All(context.Background(), &rating_result); err != nil {
		return err
	}
	rating := rating_result[0]["averageRating"].(float64)
	
	// update service rating
	svcp_collection := db.Client.Database("petpal").Collection("svcp")
	temp , err := svcp_collection.Aggregate(context.Background(), bson.A{
		bson.M{"$unwind": "$services"},
		bson.M{"$match": bson.M{"services.serviceID": service_id}},
	})
	if err != nil {
		return err
	}

	// get matches from temp
	var matches []bson.M
	if err = temp.All(context.Background(), &matches); err != nil {
		return err
	}
	
	// update service rating
	svcp_id := matches[0]["SVCPID"]
	filter = bson.M{"SVCPID": svcp_id}
	update = bson.M{"$set": bson.M{"averageRating": rating}}
	_, err = svcp_collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return nil
}
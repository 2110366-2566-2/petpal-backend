package chathistory

import (
	"context"
	"petpal-backend/src/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetHistory(db *models.MongoDB, roomId string, page int64, per int64) (*models.Chat, error) {
	collection := db.Collection("chat")
	// find chat by id

	var chat models.Chat = models.Chat{}
	filter := bson.D{{Key: "chatId", Value: roomId}}
	opts := options.FindOne().SetProjection(bson.D{{
		Key: "messages", Value: bson.D{{
			Key: "$slice", Value: []int64{(page - 1) * per, per},
		}},
	}})
	err := collection.FindOne(context.Background(), filter, opts).Decode(&chat)
	if err != nil {
		return nil, err
	}
	return &chat, nil
}

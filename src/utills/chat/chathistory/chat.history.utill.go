package chathistory

import (
	"context"
	"errors"
	"petpal-backend/src/models"
	"time"

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

func AddTextMessage(db *models.MongoDB, roomId string, content string, senderId string) error {
	// find chat by id
	chat, err := GetHistory(db, roomId, 1, 1)
	if err != nil {
		return err
	}

	collection := db.Collection("chat")

	var sender int
	if chat.User0ID == senderId {
		sender = 0
	} else if chat.User1ID == senderId {
		sender = 1
	} else {
		return errors.New("sender not in chat")
	}

	message := models.Message{
		MessageType: "text",
		Timestamp:   time.Now(),
		Content:     content,
		Sender:      sender,
	}

	filter := bson.D{{Key: "chatId", Value: roomId}}
	update := bson.D{{Key: "$push", Value: bson.D{{
		Key: "messages", Value: message,
	}}}}
	_, err = collection.UpdateOne(context.Background(), filter, update)

	if err != nil {
		return err
	}

	return err
}

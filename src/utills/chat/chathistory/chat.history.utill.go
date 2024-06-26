package chathistory

import (
	"context"
	"errors"
	"petpal-backend/src/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetChatHistory(db *models.MongoDB, roomId string, page int64, per int64) (*models.Chat, error) {
	collection := db.Collection("chat")
	// find chat by id

	var chat models.Chat = models.Chat{}
	filter := bson.D{{Key: "roomId", Value: roomId}}
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

func CreateChatHistory(db *models.MongoDB, roomId string, user0Id string, user1Id string, user0Type string, user1Type string) (*models.Chat, error) {
	collection := db.Collection("chat")
	// find chat by id
	var chat models.Chat = models.Chat{}
	chat.User0ID = user0Id
	chat.User1ID = user1Id
	chat.User0Type = user0Type
	chat.User1Type = user1Type
	chat.RoomID = roomId
	chat.Messages = []models.Message{}
	_, err := collection.InsertOne(context.Background(), chat)
	if err != nil {
		return nil, err
	}
	// Return the inserted user
	return &chat, nil
}

func AddTextMessage(db *models.MongoDB, roomId string, content string, senderId string, senderType string) error {
	// find chat by id
	chat, err := GetChatHistory(db, roomId, 1, 1)
	if err != nil {
		return err
	}

	collection := db.Collection("chat")

	var sender int
	if chat.User0ID == senderId && chat.User0Type == senderType {
		sender = 0
	} else if chat.User1ID == senderId && chat.User1Type == senderType {
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
		Key: "messages", Value: bson.M{
			"$each":     bson.A{message},
			"$position": 0,
		},
	}}}}
	_, err = collection.UpdateOne(context.Background(), filter, update)

	if err != nil {
		return err
	}

	return err
}

func UpdateChatHistoryHandler(db *models.MongoDB, roomID string, updateChatHistory bson.M) (*models.Chat, error) {
	collection := db.Collection("chat")
	var chat models.Chat
	filter := bson.D{{Key: "roomId", Value: roomID}}
	err := collection.FindOne(context.Background(), filter).Decode(&chat)
	if err != nil {
		return nil, err
	}
	update := bson.D{{Key: "$set", Value: updateChatHistory}}
	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return nil, err
	}
	updatedChat := models.Chat{}
	err = collection.FindOne(context.Background(), filter).Decode(&updatedChat)
	if err != nil {
		return nil, err
	}
	return &updatedChat, nil
}

func GetChatsById(db *models.MongoDB, id string, page int64, per int64, userType string) ([]models.Chat, error) {
	if userType != "user" && userType != "svcp" && userType != "admin" {
		return nil, errors.New("invalid user type")
	}

	// get collection
	collection := db.Collection("chat")

	// find chat by id
	filter := bson.D{{Key: "$or", Value: bson.A{
		bson.M{"user0Id": id, "user0type": userType},
		bson.M{"user1Id": id, "user1type": userType},
	}}}

	opts := options.Find().SetProjection(bson.D{{
		Key: "messages", Value: bson.D{{
			Key: "$slice", Value: []int64{0, 1},
		}},
	}}).SetSkip((page - 1) * per).SetLimit(per)

	cursor, err := collection.Find(context.Background(), filter, opts)
	if err != nil {
		return nil, err
	}

	var chats []models.Chat = []models.Chat{}
	if err := cursor.All(context.Background(), &chats); err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	return chats, err
}

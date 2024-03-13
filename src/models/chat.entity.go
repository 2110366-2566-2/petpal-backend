package models

import "time"

type Chat struct {
	RoomID      string    `json:"roomID" bson:"roomId"`
	User0ID     string    `json:"user0ID" bson:"user0Id"`
	User1ID     string    `json:"user1ID" bson:"user1Id"`
	User0Type   string    `json:"user0Type" bson:"user0type"`
	User1Type   string    `json:"user1Type" bson:"user1type"`
	DateCreated time.Time `json:"dateCreated" bson:"dateCreated"`
	Messages    []Message `json:"messages" bson:"messages"` // This is usually only part of the chat, not all
}

type Message struct {
	MessageType string    `json:"messageType" bson:"messageType"`
	Timestamp   time.Time `json:"timestamp" bson:"timestamp"`
	Content     string    `json:"message" bson:"message"`
	Sender      int       `json:"sender" bson:"sender"`
}

func (c *Chat) getUserIdAndType(userNum int) (string, string) {
	if userNum == 0 {
		return c.User0ID, c.User0Type
	}
	return c.User1ID, c.User1Type
}

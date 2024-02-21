package chat

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	Connection *websocket.Conn
	ID         string `json:"id"`
	RoomID     string `json:"roomId"`
	Username   string `json:"username"`
	Role       string `json:"role"`
	Message    chan *Message
}

type Message struct {
	Content     string `json:"content"`
	RoomID      string `json:"roomId"`
	Username    string `json:"username"`
	Role        string `json:"role"`
	MessageType string `json:messageType`
}

type MessageType string

const (
	Text  MessageType = "text"
	Emoji MessageType = "emoji"
	Image MessageType = "image"
	Video MessageType = "video"
)

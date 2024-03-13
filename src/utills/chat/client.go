package chat

import (
	"errors"
	"log"
	"time"

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
	Content     string    `json:"content"`
	RoomID      string    `json:"roomId"`
	Username    string    `json:"username"`
	Role        string    `json:"role"`
	MessageType string    `json:messageType`
	TimeStamp   time.Time `json:timeStamp`
}

type MessageType string

const (
	Text  MessageType = "text"
	Emoji MessageType = "emoji"
	Image MessageType = "image"
	Video MessageType = "video"
)

func (c *Client) writeMessage() error {
	defer func() {
		c.Connection.Close()
	}()

	for {
		message, ok := <-c.Message
		if !ok {
			return errors.New("channel closed unexpectedly")
		}
		if err := c.Connection.WriteJSON(message); err != nil {
			return err
		}
		// if message content is empty or too long
		/*
			Move this message error check into frontend due to easier implementations

				if message.MessageType == string(Text) && (len(message.Content) > 500 || len(message.Content) == 0) {
					return errors.New("Cannot send empty messsage or too long message")
				}*/
	}
}

func (c *Client) readMessage(h *Hub) error {
	defer func() {
		h.Unregister <- c
		c.Connection.Close()
	}()

	for {
		_, m, err := c.Connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			return err
		}
		msg := &Message{
			Content:     string(m),
			RoomID:      c.RoomID,
			Username:    c.Username,
			Role:        c.Role,
			MessageType: string(Text),
			TimeStamp:   time.Now(),
		}

		/* Add a message to chat history
		Implementing soon

		Chat client-server concept using web-socket

		1. Client authenticate the user from browser cookie and
		2. Client find the unique room id
		3. Send the roomid and clientid + clientUsername through the query paremeter to connect the socket server
		4. Load the chat history
		5. Close the socket connection when user logout or quit chat

		*/
		h.Broadcast <- msg
	}
}

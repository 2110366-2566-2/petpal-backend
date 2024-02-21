package chat

import "time"

type ChatRoom struct {
	ID      string             `json:"roomid"`
	Clients map[string]*Client `json:"clients"`
}

type Hub struct {
	Rooms      map[string]*ChatRoom
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *Message
}

func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[string]*ChatRoom),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *Message, 5),
	}
}
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			if _, roomExist := h.Rooms[client.RoomID]; roomExist {
				room := h.Rooms[client.RoomID]

				if _, clientExist := room.Clients[client.ID]; !clientExist {
					room.Clients[client.ID] = client
				}
			}
		case client := <-h.Unregister:
			if _, roomExist := h.Rooms[client.RoomID]; roomExist {
				if len(h.Rooms[client.RoomID].Clients) != 0 {
					h.Broadcast <- &Message{
						Content:   "user left the chat",
						RoomID:    client.RoomID,
						Username:  client.Username,
						TimeStamp: time.Now(),
					}
				}

				delete(h.Rooms[client.RoomID].Clients, client.ID)
				close(client.Message)

			}
		case msg := <-h.Broadcast:
			if _, ok := h.Rooms[msg.RoomID]; ok {

				for _, cl := range h.Rooms[msg.RoomID].Clients {
					cl.Message <- msg
				}
			}

		}
	}
}

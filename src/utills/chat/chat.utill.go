package chat

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type CreateChatRoomReq struct {
	ID string `json:"id"`
}

func CreateChatRoom(c *gin.Context, h *Hub) {
	var req CreateChatRoomReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.Rooms[req.ID] = &ChatRoom{
		ID:      req.ID,
		Clients: make(map[string]*Client),
	}

	c.JSON(http.StatusOK, req)
}

func JoinChatRoom(c *gin.Context, h *Hub) {
	connection, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	roomID := c.Param("roomId")
	clientID := c.Query("clientId")
	username := c.Query("username")
	userrole := c.Query("role")

	// Check if the roomid is existed or not ?
	if _, ok := h.Rooms[roomID]; !ok {
		// if not create a room here
		h.Rooms[roomID] = &ChatRoom{
			ID:      roomID,
			Clients: make(map[string]*Client),
		}
	}
	client := &Client{
		Connection: connection,
		ID:         clientID,
		RoomID:     roomID,
		Username:   username,
		Role:       userrole,
		Message:    make(chan *Message, 10),
	}
	//fmt.Printf("%+v\n", client)
	msg := &Message{
		Content:     "A new user has joined the room",
		RoomID:      roomID,
		UserID:      clientID,
		Username:    username,
		Role:        userrole,
		MessageType: string(Text),
		TimeStamp:   time.Now(),
	}
	h.Register <- client
	h.Broadcast <- msg

	go func() {
		err := client.writeMessage()
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}()

	err = client.readMessage(h)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
}

type ChatRoomRes struct {
	ID string `json:"id"`
}

func GetChatRooms(c *gin.Context, h *Hub) {
	rooms := make([]ChatRoomRes, 0)

	for _, r := range h.Rooms {
		rooms = append(rooms, ChatRoomRes{
			ID: r.ID,
		})
	}

	c.JSON(http.StatusOK, rooms)
}

type ClientRes struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

func GetClients(c *gin.Context, h *Hub) {
	var clients []ClientRes
	roomId := c.Param("roomId")

	if _, ok := h.Rooms[roomId]; !ok {
		clients = make([]ClientRes, 0)
		c.JSON(http.StatusOK, clients)
	}

	for _, c := range h.Rooms[roomId].Clients {
		clients = append(clients, ClientRes{
			ID:       c.ID,
			Username: c.Username,
		})
	}

	c.JSON(http.StatusOK, clients)
}

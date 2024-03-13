package routes

import (
	"petpal-backend/src/utills/chat"

	"github.com/gin-gonic/gin"
)

func ChatRoutes(r *gin.Engine, h *chat.Hub) {
	chatGroup := r.Group("/chat")
	chatGroup.POST("/createRoom", func(c *gin.Context) {
		chat.CreateChatRoom(c, h)
	})
	chatGroup.GET("/joinRoom/:roomId", func(c *gin.Context) {
		chat.JoinChatRoom(c, h)
	})
	chatGroup.GET("/getRooms", func(c *gin.Context) {
		chat.GetChatRooms(c, h)
	})
	chatGroup.GET("/getClients/:roomId", func(c *gin.Context) {
		chat.GetClients(c, h)
	})
}

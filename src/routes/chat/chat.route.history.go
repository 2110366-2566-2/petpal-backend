package routes

import (
	controllers "petpal-backend/src/controllers/chathistory"
	"petpal-backend/src/models"

	"github.com/gin-gonic/gin"
)

func ChatHistoryRoutes(r *gin.RouterGroup) {
	chatHistoryGroup := r.Group("")
	{
		chatHistoryGroup.GET("/history/:roomId", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllers.GetChatHistoryHandler(c, db)
		})
	}
}

package routes

import (
	"github.com/gin-gonic/gin"
	"petpal-backend/src/controllers/service"
	"petpal-backend/src/models"
)

func ServiceFeedbackRoutes(r *gin.RouterGroup) {
	serviceFeedbackGroup := r.Group("/feedback")

	serviceFeedbackGroup.POST("/:id", func(c *gin.Context, ) {
		db := c.MustGet("db").(*models.MongoDB)
		id := c.Param("id")
		controllers.CreateFeedbackHandler(c, db, id)
	})
}

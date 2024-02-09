package routes

import (
	"petpal-backend/src/controllers"
	"petpal-backend/src/models"

	"github.com/gin-gonic/gin"
)

func ServiceProviderRoutes(r *gin.Engine) {
	userGroup := r.Group("/svcp")
	{
		// send user profile image (Form Fields : username, profileImage)
		userGroup.POST("/uploadProfileImage", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllers.UploadImageHandler(c, "svcp", db)
		})

	}
}

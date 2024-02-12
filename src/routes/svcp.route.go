package routes

import (
	"petpal-backend/src/controllers"
	"petpal-backend/src/models"

	"github.com/gin-gonic/gin"
)

func ServiceProviderRoutes(r *gin.Engine) {
	svcpGroup := r.Group("/svcp")
	{

		//waring string in db shoudn't have '/' in string or api in input
		// update user profile image (Form Fields : username, profileImage)
		svcpGroup.POST("/uploadProfileImage", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllers.UploadImageHandler(c, "svcp", db)
		})

		// get svcp profile image (only one image)
		svcpGroup.GET("/ProfileImage/:username", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllers.GetProfileImageHandler(c, "svcp", db)
		})
	}
}

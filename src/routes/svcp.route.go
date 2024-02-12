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
		// update svcp profile image (Form Fields : email:content, profileImage:content)
		svcpGroup.POST("/uploadProfileImage", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllers.UploadImageHandler(c, "svcp", db)
		})

		//waring string in db shoudn't have '/' in string or api in input
		// get svcp profile image  (Form Fields : email:content) (only one image)
		svcpGroup.POST("/profileImage", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllers.GetProfileImageHandler(c, "svcp", db)
		})
	}
}

package routes

import (
	"petpal-backend/src/controllers"
	"petpal-backend/src/models"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine) {
	userGroup := r.Group("/user")
	{
		// userGroup.GET("/", getUserList)
		// userGroup.GET("/:id", getUserByID)
		userGroup.POST("/", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllers.CreateUserHandler(c.Writer, c.Request, db)
		})
		// userGroup.PUT("/:id", updateUser)
		// userGroup.DELETE("/:id", deleteUser)
		userGroup.POST("/setDefaultBankAccount", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllers.SetDefaultBankAccountHandler(c.Writer, c.Request, db)
		})

		//waring string in db shoudn't have '/' in string or api in input
		// update user profile image (Form Fields : email:content, profileImage:content)
		userGroup.POST("/uploadProfileImage", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllers.UploadImageHandler(c, "user", db)
		})

		//waring string in db shoudn't have '/' in string or api in input
		// get user profile image  (Form Fields : email:content) (only one image)
		userGroup.GET("/profileImage", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllers.GetProfileImageHandler(c, "user", db)
		})

	}
}

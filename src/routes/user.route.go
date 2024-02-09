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

		// send user profile image (Form Fields : username, profileImage)
		userGroup.POST("/uploadProfileImage", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllers.UploadImageHandler(c, "user", db)
		})

		// get user profile image (only one image)
		userGroup.GET("/ProfileImage/:username", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllers.GetProfileImageHandler(c, "user", db)
		})

	}
}

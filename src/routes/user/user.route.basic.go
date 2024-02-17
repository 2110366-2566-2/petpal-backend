package routes

import (
	controllers "petpal-backend/src/controllers/user"
	"petpal-backend/src/models"

	"github.com/gin-gonic/gin"
)

func UserBaseRoutes(r *gin.RouterGroup) {
	userGroup := r.Group("")
	userGroup.GET("/", func(c *gin.Context) {
		db := c.MustGet("db").(*models.MongoDB)
		controllers.GetUsersHandler(c.Writer, c.Request, db)
	})
	userGroup.GET("/:id", func(c *gin.Context) {
		db := c.MustGet("db").(*models.MongoDB)
		controllers.GetUserByIDHandler(c.Writer, c.Request, db, c.Param("id"))
	})
	userGroup.PUT("/", func(c *gin.Context) {
		db := c.MustGet("db").(*models.MongoDB)
		controllers.UpdateUserHandler(c, db)
	})

	userGroup.POST("/set-default-bank-account", func(c *gin.Context) {
		db := c.MustGet("db").(*models.MongoDB)
		controllers.SetDefaultBankAccountHandler(c.Writer, c.Request, db)
	})

	userGroup.DELETE("/delete-bank-account", func(c *gin.Context) {
		db := c.MustGet("db").(*models.MongoDB)
		controllers.DeleteBankAccountHandler(c.Writer, c.Request, db)
	})

	// update user profile image (Form Fields : email:content, profileImage:content)
	userGroup.POST("/uploadProfileImage", func(c *gin.Context) {
		db := c.MustGet("db").(*models.MongoDB)
		controllers.UploadImageHandler(c, "user", db)
	})

	// get user profile image  (Form Fields : email:content) (only one image)
	userGroup.POST("/profileImage", func(c *gin.Context) {
		db := c.MustGet("db").(*models.MongoDB)
		controllers.GetProfileImageHandler(c, "user", db)
	})
}
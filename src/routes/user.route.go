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
		userGroup.POST("/deleteBankAccount", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllers.DeleteBankAccountHandler(c.Writer, c.Request, db)
		})
	}
}

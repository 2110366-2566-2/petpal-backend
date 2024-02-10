package routes

import (
	"petpal-backend/src/controllers"
	"petpal-backend/src/models"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine) {
	userGroup := r.Group("/user")
	{
		userGroup.GET("/", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllers.GetUsersHandler(c.Writer, c.Request, db)
		})

		userGroup.GET("/:id", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllers.GetUserByIDHandler(c.Writer, c.Request, db, c.Param("id"))
		})

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
	}
}

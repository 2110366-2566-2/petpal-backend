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
	}
}

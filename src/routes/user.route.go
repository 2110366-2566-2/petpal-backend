package routes

import (
	controllers "petpal-backend/src/controllers/user"
	"petpal-backend/src/models"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine) {
	userGroup := r.Group("/user")
	{
		// userGroup.GET("/", getUserList)
		// userGroup.GET("/:id", getUserByID)
		userGroup.POST("/register", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllers.RegisterUserHandler(c, db)
		})

		userGroup.POST("/login", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllers.LoginUserHandler(c, db)
		})
		userGroup.POST("/logout", func(c *gin.Context) {
			controllers.LogoutUserHandler(c)
		})

		userGroup.GET("/me", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllers.CurrentUserHandler(c, db)
		})
		// userGroup.PUT("/:id", updateUser)
		// userGroup.DELETE("/:id", deleteUser)
		userGroup.POST("/setDefaultBankAccount", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllers.SetDefaultBankAccountHandler(c.Writer, c.Request, db)
		})
	}
}

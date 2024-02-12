package routes

import (
	controllers "petpal-backend/src/controllers/user"
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

		userGroup.PUT("/:id", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllers.UpdateUserHandler(c, db)
		})

		// userGroup.DELETE("/:id", deleteUser)
		userGroup.POST("/pets", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllers.GetUserPetsHandler(c, db)
		})
		userGroup.PUT("/:id/pets", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllers.AddUserPetHandler(c, db)
		})

		userGroup.POST("/set-default-bank-account", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllers.SetDefaultBankAccountHandler(c.Writer, c.Request, db)
		})
		userGroup.DELETE("/delete-bank-account", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllers.DeleteBankAccountHandler(c.Writer, c.Request, db)
		})
	}
}

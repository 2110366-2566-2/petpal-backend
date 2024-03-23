package routes

import (
	controllers "petpal-backend/src/controllers"
	"petpal-backend/src/models"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {
	authGroup := r.Group("")
	authGroup.GET("/current-entity", func(c *gin.Context) {
		db := c.MustGet("db").(*models.MongoDB)
		controllers.GetCurrentEntityHandler(c, db)
	})
	authGroup.POST("/register-svcp", func(c *gin.Context) {
		db := c.MustGet("db").(*models.MongoDB)
		controllers.RegisterSVCPHandler(c, db)
	})
	authGroup.POST("/register-user", func(c *gin.Context) {
		db := c.MustGet("db").(*models.MongoDB)
		controllers.RegisterUserHandler(c, db)
	})
	authGroup.POST("/register-admin", func(c *gin.Context) {
		db := c.MustGet("db").(*models.MongoDB)
		controllers.RegisterAdminHandler(c, db)
	})
	authGroup.POST("/login", func(c *gin.Context) {
		db := c.MustGet("db").(*models.MongoDB)
		controllers.LoginHandler(c, db)
	})
	authGroup.POST("/logout", func(c *gin.Context) {
		controllers.LogoutHandler(c)
	})
	authGroup.POST("/change-password", func(c *gin.Context) {
		db := c.MustGet("db").(*models.MongoDB)
		controllers.ChangePasswordHandler(c, db)
	})
}

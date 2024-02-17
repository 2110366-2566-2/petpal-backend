package routes

import (
	controllers "petpal-backend/src/controllers/user"
	"petpal-backend/src/models"

	"github.com/gin-gonic/gin"
)

func UserAuthRoutes(r *gin.RouterGroup) {
	authGroup := r.Group("")
	authGroup.POST("/register", func(c *gin.Context) {
		db := c.MustGet("db").(*models.MongoDB)
		controllers.RegisterUserHandler(c, db)
	})
	authGroup.POST("/changePassword", func(c *gin.Context) {
		db := c.MustGet("db").(*models.MongoDB)
		controllers.ChangePassword(c.Writer, c.Request, db)
	})
	authGroup.POST("/login", func(c *gin.Context) {
		db := c.MustGet("db").(*models.MongoDB)
		controllers.LoginUserHandler(c, db)
	})
	authGroup.POST("/logout", func(c *gin.Context) {
		controllers.LogoutUserHandler(c)
	})
}

package routes

import (
	"petpal-backend/src/controllers"

	"github.com/gin-gonic/gin"
)

func ExampleRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.GET("/", controllers.ExampleController())

	return router
}
func ExampleRoutes(r *gin.Engine) {
	userController := controllers.UserController{}
	userGroup := r.Group("/example")
	{
		// userGroup.GET("/", getUserList)
		// userGroup.GET("/:id", getUserByID)
		userGroup.POST("/", func(c *gin.Context) {
			userController.CreateUserHandler(c.Writer, c.Request)
		})
		// userGroup.PUT("/:id", updateUser)
		// userGroup.DELETE("/:id", deleteUser)
	}
}

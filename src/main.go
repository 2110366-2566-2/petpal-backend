package main

import (
	"net/http"
	"petpal-backend/src/configs"
	"petpal-backend/src/controllers"
	"petpal-backend/src/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize Gin router
	r := gin.Default()

	port := configs.GetPort()
	// basic route
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello World!",
		})
	})
	// add additional router
	r.GET("/example", controllers.ExampleController())
	routes.UserRoutes(r)

	r.Run(":" + port)
}

package main

import (
	"net/http"
	"petpal-backend/src/configs"
	"petpal-backend/src/controllers"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	port := configs.GetPort()
	// basic route
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello World!",
		})
	})
	// add additional router
	router.GET("/example", controllers.ExampleController())
	router.GET("/test_mongo", controllers.Testmongo())
	router.Run(":" + port)
}

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

package routes

import (
	"petpal-backend/src/controllers"

	"github.com/gin-gonic/gin"
)

func ExampleRoutes(r *gin.Engine) {
	exampleGroup := r.Group("/example")
	{
		exampleGroup.GET("/", controllers.ExampleController())
	}
}

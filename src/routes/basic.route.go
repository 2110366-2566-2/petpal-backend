package routes

import (
	"net/http"
	"petpal-backend/src/controllers"

	"github.com/gin-gonic/gin"
)

func BasicRoutes(r *gin.Engine) {
	basicGroup := r.Group("/")
	{
		basicGroup.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Hello World!",
			})
		})
		basicGroup.GET("/test_mongo", controllers.Testmongo)
	}
}

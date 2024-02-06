package main

import (
	"fmt"
	"net/http"
	"petpal-backend/src/configs"
	"petpal-backend/src/controllers"
	"petpal-backend/src/models"
	"petpal-backend/src/routes"

	"github.com/gin-gonic/gin"
)

// Middleware to inject database connection into Gin context
func DatabaseMiddleware(db *models.MongoDB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	}
}

func main() {
	// Initialize Gin router
	r := gin.Default()

	db, err := models.NewMongoDB()
	if err != nil {
		fmt.Println("Have you ever recite Namo 3 times to praise the golden-armored warrior? That so importance na :", err)
	}

	// init database to inject in gin.context
	r.Use(DatabaseMiddleware(db))

	port := configs.GetPort()
	// basic route
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello World!",
		})
	})
	r.GET("/test_mongo", controllers.Testmongo)

	// add additional router
	routes.UserRoutes(r)
	routes.ExampleRoutes(r)

	r.Run(":" + port)
}

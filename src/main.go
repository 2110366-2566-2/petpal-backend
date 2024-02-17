package main

import (
	"fmt"
	"petpal-backend/src/configs"
	"petpal-backend/src/models"
	"petpal-backend/src/routes"
	"petpal-backend/src/utills"

	"github.com/gin-gonic/gin"
)

// Middleware to inject database connection into Gin context
func DatabaseMiddleware(db *models.MongoDB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	}
}

func InitGinRouter() *gin.Engine {
	// Initialize Gin router
	r := gin.Default()

	db, err := utills.NewMongoDB()
	if err != nil {
		fmt.Println("Have you ever recite Namo 3 times to praise the golden-armored warrior? That so importance na :", err)
	}

	// init database to inject in gin.context
	// r.Use(DatabaseMiddleware(db))
	r.Use(DatabaseMiddleware(db))
	return r
}

func main() {
	// Initialize Gin router
	r := InitGinRouter()

	port := configs.GetPort()

	// add router
	routes.UserRoutes(r)
	routes.SVCPRoutes(r)

	r.Run("localhost:" + port)
}

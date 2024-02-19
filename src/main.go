package main

import (
	"fmt"
	"petpal-backend/src/configs"
	"petpal-backend/src/models"
	"petpal-backend/src/routes"
	user_route "petpal-backend/src/routes/user"
	"petpal-backend/src/utills"

	"github.com/gin-gonic/gin"

	"petpal-backend/src/docs"

	"github.com/gin-contrib/cors"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	// set cors
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}))

	port := configs.GetPort()

	// add router
	user_route.UserRoutes(r)
	routes.SVCPRoutes(r)

	// Swagger
	docs.SwaggerInfo.Title = "PetPal API"
	docs.SwaggerInfo.Description = "This is a simple API for PetPal project"
	docs.SwaggerInfo.Version = "1"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/"

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run("localhost:" + port)
}

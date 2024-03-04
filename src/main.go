package main

import (
	"fmt"
	"petpal-backend/src/configs"
	"petpal-backend/src/models"
	"petpal-backend/src/routes"
	chat_route "petpal-backend/src/routes/chat"
	service_route "petpal-backend/src/routes/service"
	user_route "petpal-backend/src/routes/user"
	"petpal-backend/src/utills"
	"petpal-backend/src/utills/chat"
	"time"

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

	// Initial chat websocket hub and run it
	h := chat.NewHub()
	go h.Run()
	location, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		fmt.Println("Error loading location:", err)
		return
	}
	time.Local = location

	// set cors
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"https://localhost:3000", "http://localhost:3000"},
		// AllowAllOrigins:  true,
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}))

	port := configs.GetPort()

	// add router
	user_route.UserRoutes(r)
	routes.SVCPRoutes(r)
	routes.AuthRoutes(r)
	service_route.ServiceRoutes(r)
	chat_route.ChatRoutes(r, h)

	// Swagger
	docs.SwaggerInfo.Title = "PetPal API"
	docs.SwaggerInfo.Description = "This is a simple API for PetPal project"
	docs.SwaggerInfo.Version = "1"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/"

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run("localhost:" + port)
}

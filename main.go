// hello world
package main

import "github.com/gin-gonic/gin"
import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Hello World!")
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello World!",
		})
	})
	router.Run(":8080")
}
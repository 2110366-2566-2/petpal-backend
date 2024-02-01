package controllers

import (
	"net/http"
	"petpal-backend/src/models"

	"github.com/gin-gonic/gin"
)

// example Go controller function
func Testmongo() gin.HandlerFunc {
	return func(c *gin.Context) {
		mongo, err := models.NewMongoDB()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		mongo.InitFirstDB()
		response, _ := mongo.GetFirstDB()
		c.JSON(http.StatusOK, response)
	}
}

package controllers

import (
	"net/http"
	"petpal-backend/src/utills"

	"github.com/gin-gonic/gin"
)

// example Go controller function
func ExampleController() gin.HandlerFunc {
	return func(c *gin.Context) {
		response, err := utills.ExampleUtill()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, response)
	}
}

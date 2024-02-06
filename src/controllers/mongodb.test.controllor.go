package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func Testmongo(c *gin.Context) {

	db := c.MustGet("db").(*mongo.Client)

	db.InitFirstDB()
	response, err := db.GetFirstDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response)
}

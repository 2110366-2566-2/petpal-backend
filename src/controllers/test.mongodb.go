package controllers

import (
	"net/http"
	"petpal-backend/src/models"

	"github.com/gin-gonic/gin"
)

func Testmongo(c *gin.Context) {

	db := c.MustGet("db").(*models.MongoDB)

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

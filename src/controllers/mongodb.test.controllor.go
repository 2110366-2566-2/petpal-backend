package controllers

import (
	"net/http"

	"petpal-backend/src/models"
	"petpal-backend/src/utills"

	"github.com/gin-gonic/gin"
)

func Testmongo(c *gin.Context) {

	db := c.MustGet("db").(*models.MongoDB)

	utills.AddMockDataToDB(db)
	response, err := utills.GetFirstDB(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response)
}

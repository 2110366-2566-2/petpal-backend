package controllers

import (
	"petpal-backend/src/models"

	"github.com/gin-gonic/gin"
)

// CreateBookingHandler godoc
//
// @Summary Create a feedback
// @Description Create a feedback for a service
// @Tags Service
//
// @Accept json
// @Produce json
//
// @Param id path string true "Service ID"
// @Param body body CreateFeedbackRequest true "Feedback rating and content(optional)"
//
// @Success 200 {object} models.Feedback
// @Failure 400 {object} models.BasicErrorRes
// @Failure 500 {object} models.BasicErrorRes
//
// @Router /service/booking/create [put]
func CreateBookingHandler(c *gin.Context, db *models.MongoDB) {
	// create booking
}

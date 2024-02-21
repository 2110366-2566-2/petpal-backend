package controllers

import (
	"github.com/gin-gonic/gin"
	"petpal-backend/src/models"
)

// CreateFeedbackHandler godoc
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
// @Router /service/feedback/{id} [post]
func CreateFeedbackHandler(c *gin.Context, db *models.MongoDB, service_id string) {

	// recalculate service rating and update service rating
}
type CreateFeedbackRequest struct {
	Content string  `json:"content" bson:"content"`
	Rating  float32 `json:"rating" bson:"rating"`
}


// GetFeedbackHandler godoc
//
// @Summary Get feedbacks
// @Description Get feedbacks for a service
// @Tags Service
//
// @Accept json
// @Produce json
//
// @Param id path string true "Service ID"
//
// @Success 200 {object} []models.Feedback
// @Failure 400 {object} models.BasicErrorRes
// @Failure 500 {object} models.BasicErrorRes
//
// @Router /service/feedback/{id} [get]
func GetFeedbackHandler(c *gin.Context, db *models.MongoDB, service_id string) {
	
	// get feedbacks for a service
}
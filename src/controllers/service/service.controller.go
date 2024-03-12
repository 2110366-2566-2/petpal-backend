package controllers

import (
	"errors"
	"net/http"
	"petpal-backend/src/models"
	"petpal-backend/src/utills/auth"
	service_utills "petpal-backend/src/utills/service"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateFeedbackHandler godoc
//
// @Summary Create a feedback
// @Description Create a feedback for a service
// @Tags Service
// @Security ApiKeyAuth
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
	// get body and create feedback
	var req CreateFeedbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, models.BasicErrorRes{Error: err.Error()})
		return
	}
	var feedback models.Feedback
	feedback.FeedbackID = primitive.NewObjectID().Hex()
	feedback.Content = req.Content
	feedback.Rating = req.Rating

	// get current user
	user, err := _authenticate(c, db)
	if err != nil {
		return
	}

	// add feedback to service
	add_err := service_utills.UpdateFeedbackToService(db, service_id, user.ID, feedback)
	if add_err != nil {
		c.JSON(500, models.BasicErrorRes{Error: add_err.Error()})
		return
	}
	c.JSON(200, feedback)
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
// @Failure 500 {object} models.BasicErrorRes
//
// @Router /service/feedback/{id} [get]
func GetFeedbackHandler(c *gin.Context, db *models.MongoDB, service_id string) {
	// get feedbacks
	feedbacks, err := service_utills.GetFeedbacksByServiceID(db, service_id)
	if err != nil {
		c.JSON(500, models.BasicErrorRes{Error: err.Error()})
		return
	}
	c.JSON(200, feedbacks)
}

func _authenticate(c *gin.Context, db *models.MongoDB) (*models.User, error) {
	entity, err := auth.GetCurrentEntityByGinContenxt(c, db)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: "Failed to get token from Cookie plase login first, " + err.Error()})
		return nil, err
	}
	switch entity := entity.(type) {
	case *models.SVCP:
		err = errors.New("need token of type User but recives token SVCP type")
		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: err.Error()})
		return nil, nil
		// Handle user
	case *models.User:
		return entity, nil
		// Handle svcp
	}
	err = errors.New("need token of type User but wrong type")
	c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: err.Error()})
	return nil, err
}

func _authenticateSVCP(c *gin.Context, db *models.MongoDB) (*models.SVCP, error) {
	entity, err := auth.GetCurrentEntityByGinContenxt(c, db)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: "Failed to get token from Cookie plase login first, " + err.Error()})
		return nil, err
	}
	switch entity := entity.(type) {
	case *models.SVCP:
		return entity, nil
		// Handle user
	case *models.User:
		err = errors.New("need token of type SVCP but recives token User type")
		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: err.Error()})
		return nil, nil
		// Handle svcp
	}
	err = errors.New("need token of type SVCP but wrong type")
	c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: err.Error()})
	return nil, err
}

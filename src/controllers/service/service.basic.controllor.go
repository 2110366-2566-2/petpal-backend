package controllers

import (
	"net/http"
	"petpal-backend/src/models"
	service_utills "petpal-backend/src/utills/service"

	"github.com/gin-gonic/gin"
)

// CreateServicesHandler godoc
//
// @Summary Create a service
// @Description Create a new service
// @Tags Service
// @Security ApiKeyAuth
//
// @Accept json
// @Produce json
//
// @Param body body CreateService true "Service data"
//
// @Success 200 {object} models.Service
// @Failure 400 {object} models.BasicErrorRes
// @Failure 500 {object} models.BasicErrorRes
//
// @Router /service/create [post]
type CreateService struct {
	Content string  `json:"content" bson:"content"`
	Rating  float32 `json:"rating" bson:"rating"`
}

func CreateServicesHandler(c *gin.Context, db *models.MongoDB) {
	// Parse request body to get user data
	var createService CreateService
	if err := c.ShouldBindJSON(&createService); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Get current user
	user, err := _authenticate(c, db)
	if err != nil {
		return
	}
	// Create a new service
	service, err := service_utills.CreateService(db, createService, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, service)
}

// SearchServicesHandler godoc
//
// @Summary Search services
// @Description Search services based on query, location, timeslot, price range, rating
// @Tags Service
//
// @Accept json
// @Produce json
//
// @Param q query string false "Search query "
// @Param location query string false "Location"
// @Param timeslot query string false "Timeslot"
// @Param start_price_range query string false "Start price range"
// @Param end_price_range query string false "End price range"
// @Param min_rating query string false "Minimum rating"
// @Param max_rating query string false "Maximum rating"
//
// @Success 200 {object} []models.Service
// @Failure 500 {object} models.BasicErrorRes
//
// @Router /service/searching [get]

func SearchServicesHandler(c *gin.Context, db *models.MongoDB, q, location, timeslot, start_price_range, end_price_range, min_rating, max_rating string) {
	// Get services
	services, err := service_utills.SearchServices(db, q, location, timeslot, start_price_range, end_price_range, min_rating, max_rating)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, services)
}

// DuplicateServicesHandler godoc
//
// @Summary Duplicate a service
// @Description Duplicate a service
// @Tags Service
// @Security ApiKeyAuth
//
// @Accept json
// @Produce json
//
// @Param id path string true "Service ID"
//
// @Success 200 {object} models.Service
// @Failure 500 {object} models.BasicErrorRes
//
// @Router /service/duplicate/{id} [post]
func DuplicateServicesHandler(c *gin.Context, db *models.MongoDB, id string) {
	// Get services
	service, err := service_utills.DuplicateService(db, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, service)
}

// DeleteServicesHandler godoc
//
// @Summary Delete a service
// @Description Delete a service
// @Tags Service
// @Security ApiKeyAuth
//
// @Accept json
// @Produce json
//
// @Param id path string true "Service ID"
//
// @Success 200 {object} models.Service
// @Failure 500 {object} models.BasicErrorRes
//
// @Router /service/{id} [delete]
func DeleteServicesHandler(c *gin.Context, db *models.MongoDB, id string) {

	err := service_utills.DeleteService(db, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Service deleted successfully"})
}

// UpdateServicesHandler godoc
//
// @Summary Update a service
// @Description Update a service
// @Tags Service
// @Security ApiKeyAuth
//
// @Accept json
// @Produce json
//
// @Param id path string true "Service ID"
// @Param body body UpdateService true "Service data"
//
// @Success 200 {object} models.Service
// @Failure 400 {object} models.BasicErrorRes
// @Failure 500 {object} models.BasicErrorRes
//
// @Router /service/{id} [patch]

type UpdateService struct {
	Content string  `json:"content" bson:"content"`
	Rating  float32 `json:"rating" bson:"rating"`
}

func UpdateServicesHandler(c *gin.Context, db *models.MongoDB, id string) {
	// Parse request body to get user data
	var updateService UpdateService
	if err := c.ShouldBindJSON(&updateService); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Get current user
	user, err := _authenticate(c, db)
	if err != nil {
		return
	}
	// Update a service
	service, err := service_utills.UpdateService(db, updateService, id, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, service)

}

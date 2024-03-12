package controllers

import (
	"net/http"
	"petpal-backend/src/models"
	service_utills "petpal-backend/src/utills/service"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
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
// @Param body body models.CreateService true "Service data"
//
// @Success 200 {object} models.Service
// @Failure 400 {object} models.BasicErrorRes
// @Failure 500 {object} models.BasicErrorRes
//
// @Router /service/create [post]
func CreateServicesHandler(c *gin.Context, db *models.MongoDB) {
	// Parse request body to get user data
	var createService *models.CreateService
	if err := c.ShouldBindJSON(&createService); err != nil {
		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: err.Error()})
		return
	}
	// Get current user
	svcp, err := _authenticateSVCP(c, db)
	if err != nil {
		return
	}
	// Create a new service
	service, err := service_utills.AddNewServices(db, createService, svcp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BasicErrorRes{Error: err.Error()})
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
// @Param start_time query string false "start_time"
// @Param end_time query string false "end_time"
// @Param start_price_range query string false "Start price range"
// @Param end_price_range query string false "End price range"
// @Param min_rating query string false "Minimum rating"
// @Param max_rating query string false "Maximum rating"
// @Param page_number query string false "Page number"
// @Param page_size query string false "Page size"
// @Param sort_by query  false "Sort by (price, rating)"
//
// @Success 200 {object} []models.Service
// @Failure 500 {object} models.BasicErrorRes
//
// @Router /service/searching [get]

// func SearchServicesHandler(c *gin.Context, db *models.MongoDB, q, location, timeslot, start_price_range, end_price_range, min_rating, max_rating string) {
// 	// Get services
// 	services, err := service_utills.SearchServices(db, q, location, timeslot, start_price_range, end_price_range, min_rating, max_rating)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, services)
// }

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
	svcp, err := _authenticateSVCP(c, db)
	if err != nil {
		return
	}
	service, err := service_utills.DuplicateService(db, id, svcp.SVCPID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BasicErrorRes{Error: err.Error()})
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
	svcp, err := _authenticateSVCP(c, db)
	if err != nil {
		return
	}
	err = service_utills.DeleteService(db, id, svcp.SVCPID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BasicErrorRes{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.BasicRes{Message: "Service deleted successfully"})
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

func UpdateServicesHandler(c *gin.Context, db *models.MongoDB, serviceID string) {
	current_svcp, err := _authenticateSVCP(c, db)
	if err != nil {
		c.JSON(http.StatusForbidden, models.BasicErrorRes{Error: "Failed to authenticate service provider " + err.Error()})
		return
	}
	var updateService *bson.M
	if err := c.ShouldBindJSON(updateService); err != nil {
		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: "Invalid request" + err.Error()})
		return
	}

	err = service_utills.UpdateService(db, serviceID, current_svcp.SVCPID, updateService)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BasicErrorRes{Error: "Failed to update service provider"})
		return
	}
	c.JSON(http.StatusOK, models.BasicRes{Message: "Service provider updated successfully"})
}

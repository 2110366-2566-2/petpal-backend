package controllers

import (
	"net/http"
	"petpal-backend/src/models"
	"petpal-backend/src/utills/auth"
	service_utills "petpal-backend/src/utills/service"
	user_utills "petpal-backend/src/utills/user"

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
// @Description Search services based on query, Address, timeslot, price range, rating
// @Tags Service
//
// @Accept json
// @Produce json
//
// @Param body body models.SearchFilter true "Search filter"
//
// @Success 200 {object} []models.Service
// @Failure 500 {object} models.BasicErrorRes
//
// @Router /service/searching [post]
func SearchServicesHandler(c *gin.Context, db *models.MongoDB) {
	// Get services
	var id string
	var is_user bool
	var searchFilter *models.SearchFilter

	if err := c.ShouldBindJSON(&searchFilter); err != nil {
		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: err.Error()})
		return
	}
	currentEntity, err := auth.GetCurrentEntityByGinContenxt(c, db)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: "Failed to get token from Cookie plase login first, " + err.Error()})
		return
	}
	switch currentEntity := currentEntity.(type) {
	case *models.SVCP:
		id = currentEntity.SVCPID
		is_user = false
		services, err := service_utills.SearchServices(db, searchFilter, id, is_user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, services)

	case *models.User:
		id = currentEntity.ID
		is_user = true
		services, err := service_utills.SearchServices(db, searchFilter, id, is_user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Add search history to user
		err = user_utills.AddSearchHistory(db, id, *searchFilter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, services)
	}
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

// GetServicesHandler godoc
//
// @Summary Get a service
// @Description Get a service
// @Tags Service
//
// @Accept json
// @Produce json
//
// @Param id path string true "Service ID"
//
// @Success 200 {object} models.Service
// @Failure 500 {object} models.BasicErrorRes
//
// @Router /service/{id} [get]
func GetServicesHandler(c *gin.Context, db *models.MongoDB, id string) {
	service, err := service_utills.GetServiceByID(db, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BasicErrorRes{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, service)
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
// @Param body body models.Service true "Service data"
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
	var updateService bson.M
	if err := c.ShouldBindJSON(&updateService); err != nil { // Pass the address of updateService
		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: "Invalid request" + err.Error()})
		return
	}

	err = service_utills.UpdateService(db, serviceID, current_svcp.SVCPID, &updateService) // Pass the address of updateService
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BasicErrorRes{Error: "Failed to update service provider"})
		return
	}
	c.JSON(http.StatusOK, models.BasicRes{Message: "Service provider updated successfully"})
}

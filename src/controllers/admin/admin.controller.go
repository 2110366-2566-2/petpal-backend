package admin

import (
	"errors"
	"net/http"
	"petpal-backend/src/models"
	admin_utills "petpal-backend/src/utills/admin"
	"petpal-backend/src/utills/auth"
	service_utills "petpal-backend/src/utills/service"

	"github.com/gin-gonic/gin"
)

// AdminDeleteServiceHandler godoc
//
// @Summary Delete a service by SVCPID and service ID
// @Description Delete a service identified by service provider ID and service ID.
// @Tags Admin
// @Security ApiKeyAuth
//
// @Produce json
//
// @Param svcpID path string true "Service provider ID"
// @Param serviceID path string true "Service ID"
//
// @Success 200 {object} models.BasicRes
// @Failure 400 {object} models.BasicErrorRes
// @Failure 401 {object} models.BasicErrorRes
// @Failure 500 {object} models.BasicErrorRes
//
// @Router /admin/service/{svcpID}/{serviceID} [delete]
func AdminDeleteServiceHandler(c *gin.Context, db *models.MongoDB, svcpID string, serviceID string) {
	_, err := _authenticateAdmin(c, db) // admin object can be used for logging in the future
	if err != nil {
		return
	}
	err = service_utills.DeleteService(db, serviceID, svcpID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BasicErrorRes{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.BasicRes{Message: "Service deleted successfully"})
}

// AdminVerifyServiceProviderHandler godoc
//
// @Summary Change the verification status of a service provider
// @Description Change the verification status of a service provider identified by service provider ID to the given verification status.
// @Tags Admin
// @Security ApiKeyAuth
//
// @Accept json
// @Produce json
//
// @Param svcpID path string true "Service provider ID"
// @Param verify body object{verify=bool} true "SVCP verification status to be set to"
//
// @Success 200 {object} models.BasicRes
// @Failure 400 {object} models.BasicErrorRes
// @Failure 401 {object} models.BasicErrorRes
// @Failure 500 {object} models.BasicErrorRes
//
// @Router /admin/serviceproviders/verify/{svcpID} [patch]
func AdminVerifyServiceProviderHandler(c *gin.Context, db *models.MongoDB, svcpID string) {
	_, err := _authenticateAdmin(c, db) // admin object can be used for logging in the future
	if err != nil {
		return
	}

	// get request body
	var req struct {
		Verify bool `json:"verify"`
	}
	err = c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: "Invalid request body"})
		return
	}

	// update service provider verification status
	err = admin_utills.AdminUpdateSVCPVerify(db, svcpID, req.Verify)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BasicErrorRes{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.BasicRes{Message: "Service provider verification status updated successfully"})
}

func _authenticateAdmin(c *gin.Context, db *models.MongoDB) (*models.Admin, error) {
	entity, err := auth.GetCurrentEntityByGinContenxt(c, db)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.BasicErrorRes{Error: "Failed to get token from Cookie plase login first, " + err.Error()})
		return nil, err
	}
	switch entity := entity.(type) {
	case *models.Admin:
		return entity, nil
		// Handle user
	default:
		err = errors.New("current entity is not of type Admin")
		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: err.Error()})
		return nil, err
	}
}
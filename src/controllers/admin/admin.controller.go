package admin

import (
	"errors"
	"net/http"
	"petpal-backend/src/models"
	admin_utills "petpal-backend/src/utills/admin"
	"petpal-backend/src/utills/auth"
	"petpal-backend/src/utills/chat/chathistory"
	service_utills "petpal-backend/src/utills/service"
	"strconv"

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

// AdminGetChatsHandler godoc
//
// @Summary 	Get chats
// @Description Get chats of the current admin. Chats are paginated. Messages only contain the latest message.
// @Tags 		Admin
//
// @Security 	ApiKeyAuth
//
// @Produce 	json
//
// @Param 		page 	query		int		false 	"Page number of chat messages (default 1)"
// @Param 		per 	query		int		false 	"Number of chat messages per page (default 10)"
//
// @Success 	200 	{object}	[]models.Chat
// @Failure 	400 	{object}	models.BasicErrorRes
// @Failure 	401 	{object}	models.BasicErrorRes
// @Failure 	500 	{object}	models.BasicErrorRes
//
// @Router /admin/chats [get]
func AdminGetChatsHandler(c *gin.Context, db *models.MongoDB) {
	current_admin, err := _authenticateAdmin(c, db) // admin object can be used for logging in the future
	if err != nil {
		return
	}

	params := c.Request.URL.Query()

	// set default values for page and per
	if !params.Has("page") {
		params.Set("page", "1")
	}
	if !params.Has("per") {
		params.Set("per", "10")
	}

	// fetch page and per from request query
	page, err_page := strconv.ParseInt(params.Get("page"), 10, 64)
	per, err_per := strconv.ParseInt(params.Get("per"), 10, 64)
	if err_page != nil || err_per != nil {
		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: "Invalid page or per number"})
		return
	}

	chats, err := chathistory.GetChatsById(db, current_admin.AdminID, page, per, "admin")
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BasicErrorRes{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, chats)
}

// AdminGetDetailBookingHandler godoc
//
// @Summary 	Admin get a booking detail by booking id
// @Description	get a booking detail by booking id
// @Tags 		Booking admin
//
// @Accept		json
// @Produce 	json
//
// @Security    ApiKeyAuth
//
// @Param       bookingID      body    models.RequestBookingId    true    "booking id"
//
// @Success 	200 {object} models.BookkingDetailRes "get detail booking"
// @Failure 	400 {object} models.BasicErrorRes
// @Failure 	401 {object} models.BasicErrorRes
// @Failure 	403 {object} models.BasicErrorRes
// @Failure 	500 {object} models.BasicErrorRes
//
// @Router 		/service/booking/detail/admin [post]
func AdminGetDetailBookingHandler(c *gin.Context, db *models.MongoDB) {
	_, err := _authenticateAdmin(c, db)
	if err != nil {
		return
	}

	request := models.RequestBookingId{}

	//400 bad request
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: err.Error()})
		return
	}

	if request.BookingID == "" {
		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: "Missing required fields"})
		return
	}

	booking, err := service_utills.GetABookingDetail(db, request.BookingID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BasicErrorRes{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.BookkingDetailRes{Message: "get detail user booking successfully", Result: *booking})
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

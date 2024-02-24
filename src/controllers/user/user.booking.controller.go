package controllers

import (
	"net/http"
	"petpal-backend/src/models"
	utills "petpal-backend/src/utills/service"
	"time"

	"github.com/gin-gonic/gin"
)

type CreateBookingHandlerSuccess struct {
	Message string         `json:"message"`
	Result  models.Booking `json:"result"`
}

// CreateBookingHandler godoc
//
// @Summary 	Create a Booking
// @Description	User can create a booking for a service
// @Tags 		Service
//
// @Accept		json
// @Produce 	json
//
// @Security    ApiKeyAuth
//
// @Param       service      body    models.BookingCreate    true    "service chosen"
//
// @Success 	201 {object} CreateBookingHandlerSuccess
// @Failure 	400 {object} models.BasicErrorRes
// @Failure 	401 {object} models.BasicErrorRes
// @Failure 	500 {object} models.BasicErrorRes
//
// @Router 		/service/booking/create [post]
func CreateBookingHandler(c *gin.Context, db *models.MongoDB) {
	// create booking
	var request struct {
		Booking models.Booking `json:"booking"`
	}

	//400 bad request
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: err.Error()})
		return
	}

	//401 not authorized
	current_user, err := _authenticate(c, db)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.BasicErrorRes{Error: err.Error()})
		return
	}
	request.Booking.UserID = current_user.ID

	// Check for required fields
	if request.Booking.ServiceID == "" || request.Booking.TimeslotID == "" || request.Booking.SVCPID == "" {
		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: "Missing required fields"})
		return
	}

	request.Booking.BookingStatus = models.BookingPending

	request.Booking.BookingTimestamp = time.Now()

	returnBooking, err := utills.InsertBooking(db, &request.Booking)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BasicErrorRes{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, CreateBookingHandlerSuccess{Message: "Booking created successfully", Result: *returnBooking})

}

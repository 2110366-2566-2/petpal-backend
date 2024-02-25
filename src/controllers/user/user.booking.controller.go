package controllers

import (
	"net/http"
	"petpal-backend/src/models"
	utills "petpal-backend/src/utills/service"
	"time"

	"github.com/gin-gonic/gin"
)

// CreateBookingHandler godoc
//
// @Summary 	Create a Booking
// @Description	User can create a booking for a service
// @Tags 		Booking
//
// @Accept		json
// @Produce 	json
//
// @Security    ApiKeyAuth
//
// @Param       service      body   models.BookingRequest    true    "service chosen"
//
// @Success 	201 {object} models.BookingBasicRes
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
	if err := c.ShouldBindJSON(&request.Booking); err != nil {
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

	c.JSON(http.StatusCreated, models.BookingBasicRes{Message: "Booking created successfully", Result: *returnBooking})

}

// UserGetAllBookingHandler godoc
//
// @Summary 	get all user booking
// @Description	get all user booking
// @Tags 		Booking
//
// @Accept		json
// @Produce 	json
//
// @Security    ApiKeyAuth
//
// @Success 	200 {array} models.BookingWithIdArrayRes
// @Failure 	401 {object} models.BasicErrorRes
// @Failure 	500 {object} models.BasicErrorRes
//
// @Router 		/service/booking/all/user [get]
func UserGetAllBookingHandler(c *gin.Context, db *models.MongoDB) {
	//401 not authorized
	current_user, err := _authenticate(c, db)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.BasicErrorRes{Error: err.Error()})
		return
	}

	bookingsList, err := utills.GetAllBookingsByUser(db, current_user.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BasicErrorRes{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.BookingWithIdArrayRes{Message: "get all user booking successfully", Result: bookingsList})
}

// UserGetIncompleteBookingHandler godoc
//
// @Summary 	get all user incomplete booking
// @Description	get only booking with status pending, paid, comfirmed (all booking that not done yet)
// @Tags 		Booking
//
// @Accept		json
// @Produce 	json
//
// @Security    ApiKeyAuth
//
// @Success 	200 {array} models.BookingWithIdArrayRes
// @Failure 	401 {object} models.BasicErrorRes
// @Failure 	500 {object} models.BasicErrorRes
//
// @Router 		/service/booking/incoming/user [get]
func UserGetIncompleteBookingHandler(c *gin.Context, db *models.MongoDB) {
	//401 not authorized
	current_user, err := _authenticate(c, db)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.BasicErrorRes{Error: err.Error()})
		return
	}

	bookingsList, err := utills.GetAllBookingsByUser(db, current_user.ID)
	var newbookingsList []models.BookingWithId

	for _, booking := range bookingsList {
		if utills.CheckBookingIsNotdone(booking) {
			newbookingsList = append(newbookingsList, booking)
		}
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BasicErrorRes{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.BookingWithIdArrayRes{Message: "get user all booking incomplete successfully", Result: newbookingsList})

}

// UserGetHistoryBookingHandler godoc
//
// @Summary 	get all user history booking
// @Description	get only booking with status completed, cancelled, expired (all booking that done)
// @Tags 		Booking
//
// @Accept		json
// @Produce 	json
//
// @Security    ApiKeyAuth
//
// @Success 	200 {array} models.BookingWithIdArrayRes
// @Failure 	401 {object} models.BasicErrorRes
// @Failure 	500 {object} models.BasicErrorRes
//
// @Router 		/service/booking/history/user [get]
func UserGetHistoryBookingHandler(c *gin.Context, db *models.MongoDB) {
	//401 not authorized
	current_user, err := _authenticate(c, db)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.BasicErrorRes{Error: err.Error()})
		return
	}

	bookingsList, err := utills.GetAllBookingsByUser(db, current_user.ID)
	var newbookingsList []models.BookingWithId

	//get only booking with status completed, cancelled, expired
	for _, booking := range bookingsList {
		if utills.CheckBookingIsDone(booking) {
			newbookingsList = append(newbookingsList, booking)
		}
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BasicErrorRes{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.BookingWithIdArrayRes{Message: "get user all booking incomplete successfully", Result: newbookingsList})

}

// UserCancelBookingHandler godoc
//
// @Summary 	user cancel booking
// @Description	can only cancel booking with status pending, paid, comfirmed (all booking that not done yet)
// @Tags 		Booking
//
// @Accept		json
// @Produce 	json
//
// @Security    ApiKeyAuth
//
// @Param       bookingID      body    models.RequestBookingId    true    "booking id"
//
// @Success 	200 {object} models.BasicRes
// @Failure 	400 {object} models.BasicErrorRes
// @Failure 	401 {object} models.BasicErrorRes
// @Failure 	500 {object} models.BasicErrorRes
//
// @Router 		/service/booking/cancel/user [post]
func UserCancelBookingHandler(c *gin.Context, db *models.MongoDB) {
	// create booking
	var request models.RequestBookingId = models.RequestBookingId{}

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

	// Check for required fields
	if request.BookingID == "" {
		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: "Missing required fields"})
		return
	}

	//get booking for checking status
	booking, err := utills.GetBooking(db, request.BookingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BasicErrorRes{Error: err.Error()})
		return
	}

	if booking.UserID != current_user.ID {
		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: "This booking is not belong to you"})
		return
	}

	//check if booking status is pending, paid, comfirmed
	if !utills.CheckBookingIsNotdone(*booking) {
		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: "Booking already done cannot be cancelled"})
		return
	}

	// if booking.BookingStatus == models.BookingPaid {
	// 	//do something like return money to user
	// }
	// if booking.BookingStatus == models.BookingComfirmed {
	// 	//do something like return money to user all sent notification to svcp?
	// }

	_, err = utills.ChangeBookingStatus(db, request.BookingID, models.BookingCanceledUser)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BasicErrorRes{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, models.BasicRes{Message: "Booking cancelled successfully"})

}

// UserRescheduleBookingHandeler godoc
//
// @Summary 	user reschedule booking
// @Description	can only reschedule booking with status pending, paid, comfirmed (all booking that not done yet)
// @Tags 		Booking
//
// @Accept		json
// @Produce 	json
//
// @Security    ApiKeyAuth
//
// @Param       booking      body    models.RequestBookingRescheduled    true    "booking id and new timeslot id"
//
// @Success 	201 {object} models.BasicRes
// @Failure 	400 {object} models.BasicErrorRes
// @Failure 	401 {object} models.BasicErrorRes
// @Failure 	500 {object} models.BasicErrorRes
//
// @Router 		/service/booking/reschedule/user [post]
func UserRescheduleBookingHandeler(c *gin.Context, db *models.MongoDB) {
	// create booking
	var request models.BookingWithId = models.BookingWithId{}

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

	// Check for required fields
	if request.BookingID == "" || request.TimeslotID == "" {
		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: "Missing required fields"})
		return
	}

	//get booking for checking status
	booking, err := utills.GetBooking(db, request.BookingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BasicErrorRes{Error: err.Error()})
		return
	}

	//check if booking is belong to user
	if booking.UserID != current_user.ID {
		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: "This booking is not belong to you"})
		return
	}

	//check if booking status is pending, paid, comfirmed
	if !utills.CheckBookingIsNotdone(*booking) {
		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: "Booking already done cannot be cancelled"})
		return
	}

	//Change booking sheduled
	_, err = utills.ChangeBookingScheduled(db, request.BookingID, request.TimeslotID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BasicErrorRes{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, models.BasicRes{Message: "Booking rescheduled successfully"})

}

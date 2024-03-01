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
// @Success 	201 {object} models.BookingInsert	"Booking created successfully"
// @Failure 	400 {object} models.BasicErrorRes
// @Failure 	401 {object} models.BasicErrorRes
// @Failure 	500 {object} models.BasicErrorRes
//
// @Router 		/service/booking/create [post]
func CreateBookingHandler(c *gin.Context, db *models.MongoDB) {
	// create booking
	request := models.Booking{}

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
	request.UserID = current_user.ID

	// Check for required fields
	if request.ServiceID == "" || request.TimeslotID == "" || request.SVCPID == "" {
		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: "Missing required fields"})
		return
	}

	request.BookingTimestamp = time.Now()

	returnBooking, err := utills.InsertBooking(db, &request)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BasicErrorRes{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, models.BookingInsert{Message: "Booking created successfully", Result: *returnBooking})

}

// UserGetAllBookingHandler godoc
//
// @Summary 	get all user booking with filter(optional)
// @Description	json body not required if you dont want to filter result
// @Description  startAfter is filter booking that has timeslot Start Before this time
// @Description  reservationType is checking booking is "incoming" or "outgoing"
// @Description  cancelStatus ,paymentStatus ,svcpConfirmed ,svcpCompleted ,userCompleted is filter booking with status 0 == false, 1 == true, 2 == dont care(or you can unuse this filed in json body)
// @Description  if dont want to filter that field dont use that field in json body
// @Description filter is and-condition(&&)
// @Description example {}
// @Description example {"reservationType":"incoming","svcpCompleted": 1,"userCompleted": 0}
// @Tags 		Booking
//
// @Accept		json
// @Produce 	json
//
// @Security    ApiKeyAuth
//
// @Param       service      body    models.RequestBookingAll    false    "get all booking with filter(optional)"
//
// @Success 	200 {array} models.BookingWithIdArrayRes "get all user booking successfully"
// @Failure 	400 {object} models.BasicErrorRes
// @Failure 	401 {object} models.BasicErrorRes
// @Failure 	500 {object} models.BasicErrorRes
//
// @Router 		/service/booking/all/user [get]
func UserGetAllBookingHandler(c *gin.Context, db *models.MongoDB) {

	request := models.RequestBookingAll{CancelStatus: 2, PaymentStatus: 2, SvcpConfirmed: 2, SvcpCompleted: 2, UserCompleted: 2}

	//401 not authorized
	current_user, err := _authenticate(c, db)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.BasicErrorRes{Error: err.Error()})
		return
	}

	//400 bad request
	if err := c.ShouldBindJSON(&request); err != nil && err.Error() != "EOF" {
		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: err.Error()})
		return
	}

	bookingsList, err := utills.GetAllBookingsByUser(db, current_user.ID)

	// print(bookingsList[0].BookingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BasicErrorRes{Error: err.Error()})
		return
	}
	if !(request.StartAfter.IsZero() && request.ReservationType == "" && request.CancelStatus == 2 && request.PaymentStatus == 2 && request.SvcpConfirmed == 2 && request.SvcpCompleted == 2 && request.UserCompleted == 2) {
		bookingsList = utills.AllBookFilter(db, bookingsList, request)
	}

	bookingsList = utills.FillSVCPDetail(db, bookingsList)

	c.JSON(http.StatusOK, models.BookingWithIdArrayRes{Message: "get all user booking successfully", Result: bookingsList})
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
// @Success 	200 {object} models.BasicRes "Booking cancelled successfully"
// @Failure 	400 {object} models.BasicErrorRes
// @Failure 	401 {object} models.BasicErrorRes
// @Failure 	500 {object} models.BasicErrorRes
//
// @Router 		/service/booking/cancel/user [post]
// func UserCancelBookingHandler(c *gin.Context, db *models.MongoDB) {
// 	// create booking
// 	var request models.RequestBookingId = models.RequestBookingId{}

// 	//400 bad request
// 	if err := c.ShouldBindJSON(&request); err != nil {
// 		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: err.Error()})
// 		return
// 	}

// 	//401 not authorized
// 	current_user, err := _authenticate(c, db)
// 	if err != nil {
// 		c.JSON(http.StatusUnauthorized, models.BasicErrorRes{Error: err.Error()})
// 		return
// 	}

// 	// Check for required fields
// 	if request.BookingID == "" {
// 		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: "Missing required fields"})
// 		return
// 	}

// 	//get booking for checking status
// 	booking, err := utills.GetBooking(db, request.BookingID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, models.BasicErrorRes{Error: err.Error()})
// 		return
// 	}

// 	if booking.UserID != current_user.ID {
// 		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: "This booking is not belong to you"})
// 		return
// 	}

// 	//check if booking status is pending, paid, comfirmed

// 	// if booking.BookingStatus == models.BookingPaid {
// 	// 	//do something like return money to user
// 	// }
// 	// if booking.BookingStatus == models.BookingComfirmed {
// 	// 	//do something like return money to user all sent notification to svcp?
// 	// }

// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, models.BasicErrorRes{Error: err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, models.BasicRes{Message: "Booking cancelled successfully"})

// }

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
// @Success 	200 {object} models.BasicRes "Booking rescheduled successfully"
// @Failure 	400 {object} models.BasicErrorRes
// @Failure 	401 {object} models.BasicErrorRes
// @Failure 	500 {object} models.BasicErrorRes
//
// @Router 		/service/booking/reschedule/user [post]
// func UserRescheduleBookingHandeler(c *gin.Context, db *models.MongoDB) {
// 	// create booking
// 	var request models.BookingWithId = models.BookingWithId{}

// 	//400 bad request
// 	if err := c.ShouldBindJSON(&request); err != nil {
// 		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: err.Error()})
// 		return
// 	}

// 	//401 not authorized
// 	current_user, err := _authenticate(c, db)
// 	if err != nil {
// 		c.JSON(http.StatusUnauthorized, models.BasicErrorRes{Error: err.Error()})
// 		return
// 	}

// 	// Check for required fields
// 	if request.BookingID == "" || request.TimeslotID == "" {
// 		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: "Missing required fields"})
// 		return
// 	}

// 	//get booking for checking status
// 	booking, err := utills.GetBooking(db, request.BookingID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, models.BasicErrorRes{Error: err.Error()})
// 		return
// 	}

// 	//check if booking is belong to user
// 	if booking.UserID != current_user.ID {
// 		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: "This booking is not belong to you"})
// 		return
// 	}

// 	//check if booking status is pending, paid, comfirmed

// 	//Change booking sheduled
// 	_, err = utills.ChangeBookingScheduled(db, request.BookingID, request.TimeslotID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, models.BasicErrorRes{Error: err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, models.BasicRes{Message: "Booking rescheduled successfully"})

// }

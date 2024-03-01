package controllers

import (
	"net/http"
	"petpal-backend/src/models"
	utills "petpal-backend/src/utills/service"
	"time"

	"github.com/gin-gonic/gin"
)

// SVCPGetAllBookingHandler godoc
//
// @Summary 	get all svcp booking with filter(optional)
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
// @Success 	200 {array} models.BookingWithIdArrayRes "get all svcp booking successfully"
// @Failure 	400 {object} models.BasicErrorRes
// @Failure 	401 {object} models.BasicErrorRes
// @Failure 	500 {object} models.BasicErrorRes
//
// @Router 		/service/booking/all/svcp [post]
func SVCPGetAllBookingHandler(c *gin.Context, db *models.MongoDB) {

	request := models.RequestBookingAll{CancelStatus: 2, PaymentStatus: 2, SvcpConfirmed: 2, SvcpCompleted: 2, UserCompleted: 2}

	current_svcp, err := _authenticate(c, db)
	if err != nil {
		http.Error(c.Writer, "Failed to get current svcp", http.StatusInternalServerError)
		return
	}

	//400 bad request
	if err := c.ShouldBindJSON(&request); err != nil && err.Error() != "EOF" {
		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: err.Error()})
		return
	}

	bookingsList, err := utills.GetAllBookingsBySVCP(db, current_svcp.SVCPID)

	// print(bookingsList[0].BookingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BasicErrorRes{Error: err.Error()})
		return
	}
	if !(request.StartAfter.IsZero() && request.ReservationType == "" && request.CancelStatus == 2 && request.PaymentStatus == 2 && request.SvcpConfirmed == 2 && request.SvcpCompleted == 2 && request.UserCompleted == 2) {
		bookingsList = utills.AllBookFilter(db, bookingsList, request)
	}

	bookingsList = utills.FillSVCPDetail(db, bookingsList)

	c.JSON(http.StatusOK, models.BookingWithIdArrayRes{Message: "get all svcp booking successfully", Result: bookingsList})

}

// SVCPComfirmBookingHandler godoc
//
// @Summary 	svcp Comfirm booking
// @Description	can only comfirm booking that is not comfirmed by svcp and not cancelled
// @Tags 		Booking
//
// @Accept		json
// @Produce 	json
//
// @Security    ApiKeyAuth
//
// @Param       bookingID      body     models.RequestBookingId   true    "booking id"
//
// @Success 	200 {object}  models.BookingBasicRes "Booking svcp completed successfully"
// @Failure 	400 {object} models.BasicErrorRes
// @Failure 	401 {object} models.BasicErrorRes
// @Failure 	500 {object} models.BasicErrorRes
//
// @Router 		/service/booking/comfirm/svcp [post]
func SVCPComfirmBookingHandler(c *gin.Context, db *models.MongoDB) {

	request := models.RequestBookingId{}

	current_svcp, err := _authenticate(c, db)
	if err != nil {
		http.Error(c.Writer, "Failed to get current svcp", http.StatusInternalServerError)
		return
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: err.Error()})
		return
	}

	booking, err := utills.GetBooking(db, request.BookingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BasicErrorRes{Error: err.Error()})
		return
	}

	if booking.SVCPID != current_svcp.SVCPID {
		c.JSON(http.StatusForbidden, models.BasicErrorRes{Error: "You are not authorized to confirm this booking"})
		return
	}

	if booking.Status.SvcpConfirmed {
		c.JSON(http.StatusForbidden, models.BasicErrorRes{Error: "Booking already confirmed by svcp"})
		return
	}

	if booking.Cancel.CancelStatus {
		c.JSON(http.StatusForbidden, models.BasicErrorRes{Error: "Booking has been cancelled"})
		return
	}

	returnBooking, err := utills.SVCPConfirmBooking(db, request.BookingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BasicErrorRes{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.BookingBasicRes{Message: "Booking svcp completed successfully", Result: *returnBooking})
}

// SVCPCompleteBookingHandler godoc
//
// @Summary 	complete a svcp booking
// @Description	can only complete not completed booking by svcp ,not cancelled ,startime is pass
// @Description if svcp not comfirmed booking it will auto comfirm first
// @Tags 		Booking
//
// @Accept		json
// @Produce 	json
//
// @Security    ApiKeyAuth
//
// @Param       bookingID      body    models.RequestBookingId    true    "booking id"
//
// @Success 	200 {object} models.BookingBasicRes "Booking completed successfully"
// @Failure 	400 {object} models.BasicErrorRes
// @Failure 	401 {object} models.BasicErrorRes
// @Failure 	500 {object} models.BasicErrorRes
//
// @Router 		/service/booking/complete/svcp [post]
func SVCPCompleteBookingHandler(c *gin.Context, db *models.MongoDB) {

	request := models.RequestBookingId{}

	current_svcp, err := _authenticate(c, db)
	if err != nil {
		http.Error(c.Writer, "Failed to get current svcp", http.StatusInternalServerError)
		return
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: err.Error()})
		return
	}

	booking, err := utills.GetBooking(db, request.BookingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BasicErrorRes{Error: err.Error()})
		return
	}

	if booking.SVCPID != current_svcp.SVCPID {
		c.JSON(http.StatusForbidden, models.BasicErrorRes{Error: "You are not authorized to complete this booking"})
		return
	}

	if booking.Status.SvcpCompleted {
		c.JSON(http.StatusForbidden, models.BasicErrorRes{Error: "Booking already completed by svcp"})
		return
	}

	if booking.Cancel.CancelStatus {
		c.JSON(http.StatusForbidden, models.BasicErrorRes{Error: "Booking has been cancelled"})
		return
	}

	if booking.StartTime.After(time.Now()) {
		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: "This booking is not started yet"})
		return
	}

	if !booking.Status.SvcpConfirmed {
		_, err := utills.SVCPConfirmBooking(db, request.BookingID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.BasicErrorRes{Error: err.Error()})
			return
		}
	}

	returnBooking, err := utills.CompleteBooking(db, request.BookingID, "svcp")
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BasicErrorRes{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.BookingBasicRes{Message: "Booking svcp completed successfully", Result: *returnBooking})
}

// SVCPCancelBookingHandler godoc
//
// @Summary 	svcp cancel booking
// @Description	can only cancel not completed by svcp booking and not cancelled
// @Tags 		Booking
//
// @Accept		json
// @Produce 	json
//
// @Security    ApiKeyAuth
//
// @Param       bookingID      body     models.RequestCancelBooking   true    "booking id"
//
// @Success 	200 {object}  models.BookingBasicRes "Booking cancelled successfully"
// @Failure 	400 {object} models.BasicErrorRes
// @Failure 	401 {object} models.BasicErrorRes
// @Failure 	500 {object} models.BasicErrorRes
//
// @Router 		/service/booking/cancel/svcp [post]
func SVCPCancelBookingHandler(c *gin.Context, db *models.MongoDB) {

	request := models.RequestCancelBooking{}

	current_svcp, err := _authenticate(c, db)
	if err != nil {
		http.Error(c.Writer, "Failed to get current svcp", http.StatusInternalServerError)
		return
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: err.Error()})
		return
	}

	// Check for required fields
	if request.BookingID == "" {
		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: "Missing required fields"})
		return
	}

	booking, err := utills.GetBooking(db, request.BookingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BasicErrorRes{Error: err.Error()})
		return
	}

	if booking.SVCPID != current_svcp.SVCPID {
		c.JSON(http.StatusForbidden, models.BasicErrorRes{Error: "You are not authorized to cancel this booking"})
		return
	}

	if booking.Status.SvcpCompleted {
		c.JSON(http.StatusForbidden, models.BasicErrorRes{Error: "Booking already completed by svcp"})
		return
	}

	if booking.Cancel.CancelStatus {
		c.JSON(http.StatusForbidden, models.BasicErrorRes{Error: "This booking is already cancelled"})
		return
	}

	if booking.Status.PaymentStatus {
		// do something like return money to user
	}

	cancel := models.BookingCancel{CancelStatus: true, CancelTimestamp: time.Now(), CancelBy: "svcp", CancelReason: request.CancelReason}

	returnBooking, err := utills.CancelBooking(db, request.BookingID, cancel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BasicErrorRes{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.BookingBasicRes{Message: "Booking svcp cancelled successfully", Result: *returnBooking})
}

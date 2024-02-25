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
// @Tags 		Booking
//
// @Accept		json
// @Produce 	json
//
// @Security    ApiKeyAuth
//
// @Param       service      body   models.BookingRequest    true    "service chosen"
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

	c.JSON(http.StatusCreated, CreateBookingHandlerSuccess{Message: "Booking created successfully", Result: *returnBooking})

}

// func RescheduledBookingHandler(c *gin.Context, db *models.MongoDB) {
// 	// create booking
// 	var request struct {
// 		Booking models.Booking `json:"booking"`
// 	}

// 	//400 bad request
// 	if err := c.ShouldBindJSON(&request.Booking); err != nil {
// 		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: err.Error()})
// 		return
// 	}

// 	//401 not authorized
// 	current_user, err := _authenticate(c, db)
// 	if err != nil {
// 		c.JSON(http.StatusUnauthorized, models.BasicErrorRes{Error: err.Error()})
// 		return
// 	}
// 	request.Booking.UserID = current_user.ID

// 	// Check for required fields
// 	if request.Booking.ServiceID == "" || request.Booking.TimeslotID == "" || request.Booking.SVCPID == "" {
// 		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: "Missing required fields"})
// 		return
// 	}

// 	request.Booking.BookingStatus = models.BookingRescheduled

// 	request.Booking.BookingTimestamp = time.Now()

// 	returnBooking, err := utills.InsertBooking(db, &request.Booking)

// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, models.BasicErrorRes{Error: err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusCreated, CreateBookingHandlerSuccess{Message: "Booking created successfully", Result: *returnBooking})

// }

type GetBookingHandlerSuccess struct {
	Message string                 `json:"message"`
	Result  []models.BookingWithId `json:"result"`
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
// @Success 	200 {array} GetBookingHandlerSuccess
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

	c.JSON(http.StatusOK, GetBookingHandlerSuccess{Message: "get all user booking successfully", Result: bookingsList})
}

// UserGetUncompleteBookingHandler godoc
//
// @Summary 	get all user uncomplete booking
// @Description	get only booking with status pending, paid, comfirmed (all booking that not done yet)
// @Tags 		Booking
//
// @Accept		json
// @Produce 	json
//
// @Security    ApiKeyAuth
//
// @Success 	200 {array} GetBookingHandlerSuccess
// @Failure 	401 {object} models.BasicErrorRes
// @Failure 	500 {object} models.BasicErrorRes
//
// @Router 		/service/booking/uncomplete/user [get]
func UserGetUncompleteBookingHandler(c *gin.Context, db *models.MongoDB) {
	//401 not authorized
	current_user, err := _authenticate(c, db)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.BasicErrorRes{Error: err.Error()})
		return
	}

	bookingsList, err := utills.GetAllBookingsByUser(db, current_user.ID)
	var newbookingsList []models.BookingWithId

	for _, booking := range bookingsList {
		if booking.BookingStatus == models.BookingPending || booking.BookingStatus == models.BookingPaid || booking.BookingStatus == models.BookingComfirmed {
			newbookingsList = append(newbookingsList, booking)
		}
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BasicErrorRes{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, GetBookingHandlerSuccess{Message: "get user all booking uncomplete successfully", Result: newbookingsList})

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
// @Success 	200 {array} GetBookingHandlerSuccess
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

	for _, booking := range bookingsList {
		if booking.BookingStatus != models.BookingRescheduled && booking.BookingStatus != models.BookingPending && booking.BookingStatus != models.BookingPaid && booking.BookingStatus != models.BookingComfirmed {
			newbookingsList = append(newbookingsList, booking)
		}
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BasicErrorRes{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, GetBookingHandlerSuccess{Message: "get user all booking uncomplete successfully", Result: newbookingsList})

}

type requestBookingId struct {
	bookingID string `json:"bookingID"`
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
// @Param       bookingID      body    requestBookingId    true    "booking id"
//
// @Success 	200 {object} models.BasicRes
// @Failure 	400 {object} models.BasicErrorRes
// @Failure 	401 {object} models.BasicErrorRes
// @Failure 	500 {object} models.BasicErrorRes
//
// @Router 		/service/booking/cancel/user [post]
func UserCancelBookingHandler(c *gin.Context, db *models.MongoDB) {
	// create booking
	var request requestBookingId

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
	if request.bookingID == "" {
		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: "Missing required fields"})
		return
	}

	//get booking for checking status
	booking, err := utills.GetBooking(db, request.bookingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BasicErrorRes{Error: err.Error()})
		return
	}

	if booking.UserID != current_user.ID {
		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: "This booking is not belong to you"})
		return
	}

	//check if booking status is pending, paid, comfirmed
	if booking.BookingStatus != models.BookingPending && booking.BookingStatus != models.BookingPaid && booking.BookingStatus != models.BookingComfirmed {
		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: "Booking already done cannot be cancelled"})
		return
	}

	// if booking.BookingStatus == models.BookingPaid {
	// 	//do something like return money to user
	// }
	// if booking.BookingStatus == models.BookingComfirmed {
	// 	//do something like return money to user all sent notification to svcp?
	// }

	_, err = utills.ChangeBookingStatus(db, request.bookingID, models.BookingCanceledUser)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BasicErrorRes{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, models.BasicRes{Message: "Booking cancelled successfully"})

}

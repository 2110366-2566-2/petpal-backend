package controllers

import (
	"net/http"
	"petpal-backend/src/configs"
	"petpal-backend/src/models"
	payment_utills "petpal-backend/src/utills/payment"
	service_utills "petpal-backend/src/utills/service"

	"github.com/gin-gonic/gin"
)

// GetPromptpayQrHandler godoc
// @Summary Get promptpayQr from a booking
// @Description Get promptpayQr from a booking
// @Tags Service
// @Accept json
// @Produce json
// @Param requestBody body models.RequestBookingId true "Request Body"
// @Success 200 {array} models.PromptpayQr "Success"
// @Failure 400 {object} models.BasicErrorRes "Bad Request"
// @Failure 500 {object} models.BasicErrorRes "Internal Server Error"
// @Router /service/booking/promptpayqr [post]
func GetPromptpayQrHandler(c *gin.Context, db *models.MongoDB) {
	request := models.RequestBookingId{}
	//400 bad request
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: err.Error()})
		return
	}
	if request.BookingID == "" {
		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: "Missing bookingID"})
		return
	}

	booking, err := service_utills.GetABookingDetail(db, request.BookingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BasicErrorRes{Error: err.Error()})
		return
	}
	qr, err := payment_utills.GeneratePromptpayQr(configs.GetPetpalPhoneNumber(), int(booking.TotalBookingPrice))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BasicErrorRes{Error: err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, models.PromptpayQr{QrImage: qr})
}

// AuthorizePaymentHandler godoc
// @Summary Get promptpayQr from a booking
// @Description Get promptpayQr from a booking
// @Tags Service
// @Accept json
// @Produce json
// @Param requestBody body models.RequestBookingId true "Request Body"
// @Success 200 {array} models.PromptpayQr "Success"
// @Failure 400 {object} models.BasicErrorRes "Bad Request"
// @Failure 500 {object} models.BasicErrorRes "Internal Server Error"
// @Router /service/booking/promptpayqr [post]
func AuthorizePaymentHandler(c *gin.Context, db *models.MongoDB) {
	request := models.RequestBookingId{}
	//400 bad request
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: err.Error()})
		return
	}
	if request.BookingID == "" {
		c.JSON(http.StatusBadRequest, models.BasicErrorRes{Error: "Missing bookingID"})
		return
	}

	booking, err := service_utills.GetABookingDetail(db, request.BookingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BasicErrorRes{Error: err.Error()})
		return
	}
	qr, err := payment_utills.GeneratePromptpayQr(configs.GetPetpalPhoneNumber(), int(booking.TotalBookingPrice))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.BasicErrorRes{Error: err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, models.PromptpayQr{QrImage: qr})
}

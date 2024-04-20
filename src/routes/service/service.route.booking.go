package routes

import (
	// controllers "petpal-backend/src/controllers/user"
	// "petpal-backend/src/models"

	controllersAdmin "petpal-backend/src/controllers/admin"
	controllersService "petpal-backend/src/controllers/service"
	controllersSVCP "petpal-backend/src/controllers/serviceprovider"
	controllersUser "petpal-backend/src/controllers/user"
	"petpal-backend/src/models"

	"github.com/gin-gonic/gin"
)

func ServiceBookingRoutes(r *gin.RouterGroup) {

	bookingGroup := r.Group("/booking")
	{
		bookingGroup.POST("/create", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllersUser.CreateBookingHandler(c, db)
		})

		// Get all booking for user
		bookingGroup.POST("/all/user", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllersUser.UserGetAllBookingHandler(c, db)
		})

		// Get detail for booking ID
		bookingGroup.POST("/detail/user", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllersUser.UserGetDetailBookingHandler(c, db)
		})

		bookingGroup.PATCH("/cancel/user", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllersUser.UserCancelBookingHandler(c, db)
		})

		bookingGroup.PATCH("/reschedule/user", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllersUser.UserRescheduleBookingHandeler(c, db)
		})

		bookingGroup.PATCH("/complete/user", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllersUser.UserCompleteBookingHandler(c, db)
		})

		// svcp
		bookingGroup.POST("/all/svcp", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllersSVCP.SVCPGetAllBookingHandler(c, db)
		})
		bookingGroup.PATCH("/confirm/svcp", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllersSVCP.SVCPConfirmBookingHandler(c, db)
		})
		bookingGroup.PATCH("complete/svcp", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllersSVCP.SVCPCompleteBookingHandler(c, db)
		})
		bookingGroup.PATCH("/cancel/svcp", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllersSVCP.SVCPCancelBookingHandler(c, db)
		})

		bookingGroup.POST("/detail/svcp", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllersSVCP.SVCPGetDetailBookingHandler(c, db)
		})
		// Payment
		bookingGroup.POST("/payment/qr", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllersService.GetPromptpayQrHandler(c, db)
		})
		bookingGroup.POST("/payment/authorize", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllersService.AuthorizePaymentHandler(c, db)
		})
		bookingGroup.POST("/payment/refund", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllersService.RefundBookingHandler(c, db)
		})

		bookingGroup.POST("/detail/admin", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllersAdmin.AdminGetDetailBookingHandler(c, db)
		})
	}
}

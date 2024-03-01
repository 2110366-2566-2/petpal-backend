package routes

import (
	// controllers "petpal-backend/src/controllers/user"
	// "petpal-backend/src/models"

	controllersSVCP "petpal-backend/src/controllers/serviceprovider"
	controllers "petpal-backend/src/controllers/user"

	"petpal-backend/src/models"

	"github.com/gin-gonic/gin"
)

func ServiceBookingRoutes(r *gin.RouterGroup) {

	bookingGroup := r.Group("/booking")
	{
		bookingGroup.POST("/create", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllers.CreateBookingHandler(c, db)
		})

		bookingGroup.POST("/all/user", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllers.UserGetAllBookingHandler(c, db)
		})

		bookingGroup.POST("/detail/user", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllers.UserGetDetailBookingHandler(c, db)
		})

		bookingGroup.POST("/cancel/user", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllers.UserCancelBookingHandler(c, db)
		})

		bookingGroup.POST("/reschedule/user", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllers.UserRescheduleBookingHandeler(c, db)
		})

		bookingGroup.POST("/complete/user", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllers.UserCompleteBookingHandler(c, db)
		})

		bookingGroup.POST("/all/svcp", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllersSVCP.SVCPGetAllBookingHandler(c, db)
		})
		bookingGroup.POST("/comfirm/svcp", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllersSVCP.SVCPComfirmBookingHandler(c, db)
		})
		bookingGroup.POST("complete/svcp", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllersSVCP.SVCPCompleteBookingHandler(c, db)
		})
		bookingGroup.POST("/cancel/svcp", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllersSVCP.SVCPCancelBookingHandler(c, db)
		})

		bookingGroup.POST("/detail/svcp", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllersSVCP.SVCPGetDetailBookingHandler(c, db)
		})

	}
}

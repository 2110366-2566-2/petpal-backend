package routes

import (
	// controllers "petpal-backend/src/controllers/user"
	// "petpal-backend/src/models"

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

		// bookingGroup.GET("/incoming/user", func(c *gin.Context) {
		// 	db := c.MustGet("db").(*models.MongoDB)
		// 	controllers.UserGetIncompleteBookingHandler(c, db)
		// })

		// bookingGroup.GET("/history/user", func(c *gin.Context) {
		// 	db := c.MustGet("db").(*models.MongoDB)
		// 	controllers.UserGetHistoryBookingHandler(c, db)
		// })

		// bookingGroup.POST("/cancel/user", func(c *gin.Context) {
		// 	db := c.MustGet("db").(*models.MongoDB)
		// 	controllers.UserCancelBookingHandler(c, db)
		// })

		// bookingGroup.POST("/reschedule/user", func(c *gin.Context) {
		// 	db := c.MustGet("db").(*models.MongoDB)
		// 	controllers.UserRescheduleBookingHandeler(c, db)
		// })

		// bookingGroup.GET("/all/svcp", func(c *gin.Context) {
		// 	db := c.MustGet("db").(*models.MongoDB)
		// 	//controllers.SVCPGetAllBookingHandler(c, db)
		// })
		// bookingGroup.GET("/history/svcp", func(c *gin.Context) {
		// 	db := c.MustGet("db").(*models.MongoDB)
		// 	//controllers.SVCPGetHistoryBookingHandler(c, db)
		// })
		// bookingGroup.POST("/accept/svcp", func(c *gin.Context) {
		// 	db := c.MustGet("db").(*models.MongoDB)
		// 	//controllers.SVCPAcceptBookingHandler(c, db)
		// })
		// bookingGroup.POST("/cancel/svcp", func(c *gin.Context) {
		// 	db := c.MustGet("db").(*models.MongoDB)
		// 	//controllers.SVCPRejectBookingHandler(c, db)
		// })

	}
}

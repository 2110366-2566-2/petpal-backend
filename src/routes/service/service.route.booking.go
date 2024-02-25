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

		bookingGroup.GET("/all/user", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllers.UserGetAllBookingHandler(c, db)
		})

		bookingGroup.GET("/incomplete/user", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllers.UserGetIncompleteBookingHandler(c, db)
		})

		bookingGroup.GET("/history/user", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllers.UserGetHistoryBookingHandler(c, db)
		})

		bookingGroup.POST("/cancel/user", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllers.UserCancelBookingHandler(c, db)
		})

		bookingGroup.POST("/reschedule/user", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllers.UserRescheduleBookingHandeler(c, db)
		})

	}
}

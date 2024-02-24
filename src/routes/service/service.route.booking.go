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

		bookingGroup.GET("/uncomplete/user", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllers.UserGetUncompleteBookingHandler(c, db)
		})

		bookingGroup.GET("/history/user", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllers.UserGetHistoryBookingHandler(c, db)
		})

	}
}

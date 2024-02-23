package routes

import (
	// controllers "petpal-backend/src/controllers/user"
	// "petpal-backend/src/models"

	controllers "petpal-backend/src/controllers/service"
	"petpal-backend/src/models"

	"github.com/gin-gonic/gin"
)

func ServiceBookingRoutes(r *gin.RouterGroup) {

	bookingGroup := r.Group("/booking")
	{
		bookingGroup.PUT("/create", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllers.CreateBookingHandler(c, db)
		})

	}
}

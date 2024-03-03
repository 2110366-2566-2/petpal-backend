package routes

import (
	controllers "petpal-backend/src/controllers/services"
	"petpal-backend/src/models"

	"github.com/gin-gonic/gin"
)

func ServiceBaseRoutes(r *gin.RouterGroup) {
	baiscGroup := r.Group("")
	{
		baiscGroup.POST("/create", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllers.CreateServicesHandler(c, db)
		})
		baiscGroup.GET("/searching", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			q := c.Query("q")
			location := c.Query("location")
			timeslot := c.Query("timeslot")
			start_price_range := c.Query("start_price_range")
			end_price_range := c.Query("end_price_range")
			min_rating := c.Query("min_rating")
			max_rating := c.Query("max_rating")
			controllers.SearchServicesHandler(c, db, q, location, timeslot, start_price_range, end_price_range, min_rating, max_rating)
		})
		baiscGroup.POST("/duplicate/:id", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			id := c.Param("id")
			controllers.DuplicateServicesHandler(c, db, id)
		})
		baiscGroup.DELETE("/:id", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			id := c.Param("id")
			controllers.DeleteServicesHandler(c, db, id)
		})
		baiscGroup.PATCH("/:id", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			id := c.Param("id")
			controllers.UpdateServicesHandler(c, db, id)
		})
	}
}

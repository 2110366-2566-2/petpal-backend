package routes

import (
	controllers "petpal-backend/src/controllers/service"
	"petpal-backend/src/models"

	"github.com/gin-gonic/gin"
)

func ServiceTimeslotRoutes(r *gin.RouterGroup) {
	baiscGroup := r.Group("timeslot")
	{
		baiscGroup.POST("/create", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllers.CreateServicesHandler(c, db)
		})
		baiscGroup.POST("/searching", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllers.SearchServicesHandler(c, db)
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

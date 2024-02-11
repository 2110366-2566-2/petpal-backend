package routes

import (
	controllers "petpal-backend/src/controllers/serviceprovider"
	"petpal-backend/src/models"

	"github.com/gin-gonic/gin"
)

func SVCPRoutes(r *gin.Engine) {
	SVCPGroup := r.Group("/serviceprovider")
	{

		SVCPGroup.GET("/me", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllers.CurrentSVCPHandler(c, db)
		})
		SVCPGroup.GET("/", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllers.GetSVCPsHandler(c.Writer, c.Request, db)
		})
		SVCPGroup.GET("/:id", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllers.GetSVCPByIDHandler(c.Writer, c.Request, db, c.Param("id"))
		})
		// SVCPGroup.POST("/", createSVCP)
		SVCPGroup.PUT("/:id", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllers.UpdateSVCPHandler(c.Writer, c.Request, db, c.Param("id"))
		})
		SVCPGroup.POST("/", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllers.RegisterSVCPHandler(c, db)
		})
		// SVCPGroup.DELETE("/:id", deleteSVCP)
	}
}

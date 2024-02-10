package routes

import (
	"petpal-backend/src/controllers/serviceprovider"
	"petpal-backend/src/models"

	"github.com/gin-gonic/gin"
)

func SVCPRoutes(r *gin.Engine) {
	SVCPGroup := r.Group("/serviceprovider")
	{
		SVCPGroup.GET("/", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllers.GetSVCPsHandler(c.Writer, c.Request, db)
		})
		SVCPGroup.GET("/:id", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllers.GetSVCPByIDHandler(c.Writer, c.Request, db, c.Param("id"))
		})
		// SVCPGroup.POST("/", func(c *gin.Context) {
		// 	db := c.MustGet("db").(*models.MongoDB)
		// 	controllers.CreateUserHandler(c.Writer, c.Request, db)
		// })
		// SVCPGroup.PUT("/:id", updateUser)
		// SVCPGroup.DELETE("/:id", deleteUser)
		// SVCPGroup.POST("/setDefaultBankAccount", func(c *gin.Context) {
		// 	db := c.MustGet("db").(*models.MongoDB)
		// 	controllers.SetDefaultBankAccountHandler(c.Writer, c.Request, db)
		// })
	}
}

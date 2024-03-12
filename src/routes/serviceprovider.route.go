package routes

import (
	controllers "petpal-backend/src/controllers/serviceprovider"
	user_controllers "petpal-backend/src/controllers/user"
	"petpal-backend/src/models"

	"github.com/gin-gonic/gin"
)

func SVCPRoutes(r *gin.Engine) {
	SVCPGroup := r.Group("/serviceproviders")
	{
		SVCPGroup.GET("/", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllers.GetSVCPsHandler(c, db)
		})
		SVCPGroup.GET("/:id", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllers.GetSVCPByIDHandler(c, db, c.Param("id"))
		})
		SVCPGroup.PUT("/", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllers.UpdateSVCPHandler(c, db)
		})
		// SVCPGroup.DELETE("/:id", deleteSVCP)
		SVCPGroup.POST("/set-default-bank-account", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllers.SetDefaultBankAccountHandler(c, db)
		})
		SVCPGroup.POST("/upload-description", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllers.UploadDescriptionHandler(c, db)
		})
		SVCPGroup.POST("/add-service", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllers.AddServiceHandler(c, db)
		})
		SVCPGroup.POST("/upload-license", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllers.UploadSVCPLicenseHandler(c, db)
		})
		SVCPGroup.DELETE("/delete-bank-account", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllers.DeleteBankAccountHandler(c, db)
		})

		// update svcp profile image (Form Fields : email:content, profileImage:content)
		SVCPGroup.POST("/uploadProfileImage", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			user_controllers.UploadImageHandler(c, db)
		})
	}
}

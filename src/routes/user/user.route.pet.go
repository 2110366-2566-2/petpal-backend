package routes

import (
	controllers "petpal-backend/src/controllers/user"
	"petpal-backend/src/models"

	"github.com/gin-gonic/gin"
)

func UserPetRoutes(r *gin.RouterGroup) {

	petGroup := r.Group("/pets")
	{
		petGroup.GET("/:id", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllers.GetUserPetsByIdHandler(c, db, c.Param("id"))
		})
		// userGroup.DELETE("/:id", deleteUser)
		petGroup.GET("/me", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllers.GetCurrentUserPetsHandler(c, db)
		})
		petGroup.PUT("/", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllers.AddUserPetHandler(c, db)
		})
		petGroup.DELETE("/:idx", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllers.DeleteUserPetHandler(c, db, c.Param("idx"))
		})
		petGroup.PUT("/:idx", func(c *gin.Context) {
			db := c.MustGet("db").(*models.MongoDB)
			controllers.UpdateUserPetHandler(c, db, c.Param("idx"))
		})

	}
}

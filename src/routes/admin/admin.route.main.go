package routes

import (
	"petpal-backend/src/controllers/admin"
	"petpal-backend/src/models"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(r *gin.Engine) {
	adminGroup := r.Group("/admin")

	adminGroup.DELETE("/service/:svcpID/:serviceID", func(c *gin.Context) {
		db := c.MustGet("db").(*models.MongoDB)
		svcpID := c.Param("svcpID")
		serviceID := c.Param("serviceID")
		admin.AdminDeleteServiceHandler(c, db, svcpID, serviceID)
	})
}

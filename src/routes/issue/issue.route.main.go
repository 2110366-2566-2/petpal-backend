package routes

import (
	controllers "petpal-backend/src/controllers/issue"
	"petpal-backend/src/models"

	"github.com/gin-gonic/gin"
)

func IssueRoutes(r *gin.Engine) {
	issueGroup := r.Group("/issue")

	issueGroup.POST("/", func(c *gin.Context) {
		db := c.MustGet("db").(*models.MongoDB)
		controllers.CreateIssue(c, db)
	})
}

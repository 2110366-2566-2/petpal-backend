package routes

import (
	"github.com/gin-gonic/gin"
)

func ServiceRoutes(r *gin.Engine) {
	serviceGroup := r.Group("/service")

	ServiceBaseRoutes(serviceGroup)
	ServiceFeedbackRoutes(serviceGroup)
}

package routes

import (
	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine) {
	userGroup := r.Group("/user")
	UserBaseRoutes(userGroup)
	UserPetRoutes(userGroup)
}

package routes

import (
	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine) {
	userGroup := r.Group("/user")
	UserAuthRoutes(userGroup)
	UserBaseRoutes(userGroup)
	UserPetRoutes(userGroup)
}

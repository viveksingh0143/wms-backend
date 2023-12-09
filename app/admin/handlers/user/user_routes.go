package user

import (
	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(r *gin.RouterGroup, userHandler *Handler) {
	api := r.Group("/users")
	{
		api.DELETE("/bulk-update", userHandler.DeleteBulk)
		api.GET("", userHandler.List)
		api.POST("", userHandler.Create)
		api.GET("/:id", userHandler.Get)
		api.PUT("/:id", userHandler.Update)
		api.DELETE("/:id", userHandler.Delete)
	}
}

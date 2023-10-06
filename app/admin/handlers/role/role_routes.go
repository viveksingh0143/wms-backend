package role

import (
	"github.com/gin-gonic/gin"
)

func SetupRoleRoutes(r *gin.RouterGroup, roleHandler *Handler) {
	api := r.Group("/roles")
	{
		api.DELETE("/bulk-update", roleHandler.DeleteBulk)
		api.GET("/", roleHandler.List)
		api.POST("/", roleHandler.Create)
		api.GET("/:id", roleHandler.Get)
		api.PUT("/:id", roleHandler.Update)
		api.DELETE("/:id", roleHandler.Delete)
	}
}

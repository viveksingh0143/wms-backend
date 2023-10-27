package permission

import (
	"github.com/gin-gonic/gin"
)

func SetupPermissionRoutes(r *gin.RouterGroup, permissionHandler *Handler) {
	api := r.Group("/permissions")
	{
		api.DELETE("/bulk-update", permissionHandler.DeleteBulk)
		api.GET("", permissionHandler.List)
		api.POST("", permissionHandler.Create)
		api.GET("/:id", permissionHandler.Get)
		api.PUT("/:id", permissionHandler.Update)
		api.DELETE("/:id", permissionHandler.Delete)
	}
}

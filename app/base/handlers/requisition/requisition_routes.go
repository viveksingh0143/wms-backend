package requisition

import (
	"github.com/gin-gonic/gin"
)

func SetupRequisitionRoutes(r *gin.RouterGroup, handler *Handler) {
	api := r.Group("/requisitions")
	{
		api.DELETE("/bulk-update", handler.DeleteBulk)
		api.GET("", handler.List)
		api.POST("", handler.Create)
		api.GET("/:id", handler.Get)
		api.PUT("/:id", handler.Update)
		api.DELETE("/:id", handler.Delete)
	}
}

package inventory

import (
	"github.com/gin-gonic/gin"
)

func SetupInventoryRoutes(r *gin.RouterGroup, inventoryHandler *Handler) {
	api := r.Group("/inventories")
	{
		api.DELETE("/bulk-update", inventoryHandler.DeleteBulk)
		api.GET("", inventoryHandler.List)
		api.POST("", inventoryHandler.Create)
		api.GET("/:id", inventoryHandler.Get)
		api.PUT("/:id", inventoryHandler.Update)
		api.DELETE("/:id", inventoryHandler.Delete)
	}
}

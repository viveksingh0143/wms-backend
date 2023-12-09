package inventory

import (
	"github.com/gin-gonic/gin"
)

func SetupInventoryRoutes(r *gin.RouterGroup, inventoryHandler *Handler) {
	api := r.Group("/inventories")
	{
		api.GET("", inventoryHandler.List)
		api.GET("/:id", inventoryHandler.Get)
		api.POST("/rawmaterial/stockin", inventoryHandler.RawMaterialStockin)
		api.POST("/finishedgoods/stockin", inventoryHandler.FinishedGoodsStockin)
		api.POST("/finishedgood/stockin", inventoryHandler.FinishedGoodStockin)
		api.POST("/attach-container", inventoryHandler.AttachContainer)
	}
}

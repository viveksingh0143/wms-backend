package batchlabel

import (
	"github.com/gin-gonic/gin"
	"star-wms/app/warehouse/handlers/sticker"
)

func SetupBatchlabelRoutes(r *gin.RouterGroup, batchlabelHandler *Handler, stickerHandler *sticker.Handler) {
	api := r.Group("/batchlabels")
	{
		api.DELETE("/bulk-update", batchlabelHandler.DeleteBulk)
		api.GET("", batchlabelHandler.List)
		api.POST("", batchlabelHandler.Create)
		api.GET("/:batchlabelID", batchlabelHandler.Get)
		api.PUT("/:batchlabelID", batchlabelHandler.Update)
		api.DELETE("/:batchlabelID", batchlabelHandler.Delete)

		api.GET("/:batchlabelID/stickers", stickerHandler.List)
		api.POST("/:batchlabelID/stickers", stickerHandler.Create)
		api.GET("/:batchlabelID/stickers/:id", stickerHandler.Get)
	}
}

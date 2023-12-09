package sticker

import (
	"github.com/gin-gonic/gin"
)

func SetupStickerRoutes(r *gin.RouterGroup, sticker *Handler) {
	api := r.Group("/stickers")
	{
		api.GET("", sticker.List)
		api.POST("", sticker.Create)
		api.GET("/:id", sticker.Get)
	}
}

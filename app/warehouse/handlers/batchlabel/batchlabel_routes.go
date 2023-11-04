package batchlabel

import (
	"github.com/gin-gonic/gin"
)

func SetupBatchlabelRoutes(r *gin.RouterGroup, batchlabelHandler *Handler) {
	api := r.Group("/batchlabels")
	{
		api.DELETE("/bulk-update", batchlabelHandler.DeleteBulk)
		api.GET("", batchlabelHandler.List)
		api.POST("", batchlabelHandler.Create)
		api.GET("/:id", batchlabelHandler.Get)
		api.PUT("/:id", batchlabelHandler.Update)
		api.DELETE("/:id", batchlabelHandler.Delete)
	}
}

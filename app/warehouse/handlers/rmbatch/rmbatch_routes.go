package rmbatch

import (
	"github.com/gin-gonic/gin"
)

func SetupRMBatchRoutes(r *gin.RouterGroup, handler *Handler) {
	api := r.Group("/rmbatches")
	{
		api.GET("", handler.List)
		api.GET("/:id", handler.Get)
	}
}

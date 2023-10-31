package store

import (
	"github.com/gin-gonic/gin"
)

func SetupStoreRoutes(r *gin.RouterGroup, storeHandler *Handler) {
	api := r.Group("/stores")
	{
		api.DELETE("/bulk-update", storeHandler.DeleteBulk)
		api.GET("", storeHandler.List)
		api.POST("", storeHandler.Create)
		api.GET("/:id", storeHandler.Get)
		api.PUT("/:id", storeHandler.Update)
		api.DELETE("/:id", storeHandler.Delete)
	}
}

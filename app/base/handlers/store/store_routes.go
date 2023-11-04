package store

import (
	"github.com/gin-gonic/gin"
	"star-wms/app/base/handlers/storelocation"
)

func SetupStoreRoutes(r *gin.RouterGroup, storeHandler *Handler, storelocationHandler *storelocation.Handler) {
	api := r.Group("/stores")
	{
		api.DELETE("/bulk-update", storeHandler.DeleteBulk)
		api.GET("", storeHandler.List)
		api.POST("", storeHandler.Create)
		api.GET("/:storeID", storeHandler.Get)
		api.PUT("/:storeID", storeHandler.Update)
		api.DELETE("/:storeID", storeHandler.Delete)
		api.GET("/:storeID/locations", storelocationHandler.List)
		api.POST("/:storeID/locations", storelocationHandler.Create)
		api.GET("/:storeID/locations/:id", storelocationHandler.Get)
		api.PUT("/:storeID/locations/:id", storelocationHandler.Update)
		api.DELETE("/:storeID/locations/:id", storelocationHandler.Delete)
	}
}

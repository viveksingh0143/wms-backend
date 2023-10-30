package product

import (
	"github.com/gin-gonic/gin"
)

func SetupProductRoutes(r *gin.RouterGroup, productHandler *Handler) {
	api := r.Group("/products")
	{
		api.DELETE("/bulk-update", productHandler.DeleteBulk)
		api.GET("", productHandler.List)
		api.POST("", productHandler.Create)
		api.GET("/:id", productHandler.Get)
		api.PUT("/:id", productHandler.Update)
		api.DELETE("/:id", productHandler.Delete)
	}
}

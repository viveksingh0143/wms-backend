package category

import (
	"github.com/gin-gonic/gin"
)

func SetupCategoryRoutes(r *gin.RouterGroup, categoryHandler *Handler) {
	api := r.Group("/categories")
	{
		api.DELETE("/bulk-update", categoryHandler.DeleteBulk)
		api.GET("", categoryHandler.List)
		api.POST("", categoryHandler.Create)
		api.GET("/:id", categoryHandler.Get)
		api.PUT("/:id", categoryHandler.Update)
		api.DELETE("/:id", categoryHandler.Delete)
	}
}

package plant

import (
	"github.com/gin-gonic/gin"
)

func SetupPlantRoutes(r *gin.RouterGroup, plantHandler *Handler) {
	api := r.Group("/plants")
	{
		api.DELETE("/bulk-update", plantHandler.DeleteBulk)
		api.GET("/", plantHandler.List)
		api.POST("/", plantHandler.Create)
		api.GET("/:id", plantHandler.Get)
		api.PUT("/:id", plantHandler.Update)
		api.DELETE("/:id", plantHandler.Delete)
	}
}

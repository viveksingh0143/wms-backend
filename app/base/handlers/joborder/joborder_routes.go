package joborder

import (
	"github.com/gin-gonic/gin"
)

func SetupJoborderRoutes(r *gin.RouterGroup, joborderHandler *Handler) {
	api := r.Group("/joborders")
	{
		api.DELETE("/bulk-update", joborderHandler.DeleteBulk)
		api.GET("", joborderHandler.List)
		api.POST("", joborderHandler.Create)
		api.GET("/:id", joborderHandler.Get)
		api.PUT("/:id", joborderHandler.Update)
		api.DELETE("/:id", joborderHandler.Delete)
	}
}

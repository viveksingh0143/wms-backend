package customer

import (
	"github.com/gin-gonic/gin"
)

func SetupCustomerRoutes(r *gin.RouterGroup, customerHandler *Handler) {
	api := r.Group("/customers")
	{
		api.DELETE("/bulk-update", customerHandler.DeleteBulk)
		api.GET("", customerHandler.List)
		api.POST("", customerHandler.Create)
		api.GET("/:id", customerHandler.Get)
		api.PUT("/:id", customerHandler.Update)
		api.DELETE("/:id", customerHandler.Delete)
	}
}

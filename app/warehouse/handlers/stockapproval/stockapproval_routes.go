package stockapproval

import (
	"github.com/gin-gonic/gin"
)

func SetupStockapprovalsRoutes(r *gin.RouterGroup, handler *Handler) {
	api := r.Group("/stockapprovals")
	{
		api.GET("", handler.List)
		api.DELETE("/:id", handler.Approve)
		api.DELETE("/bulk-update", handler.ApproveBulk)
	}
}

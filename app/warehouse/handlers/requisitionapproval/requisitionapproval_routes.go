package requisitionapproval

import (
	"github.com/gin-gonic/gin"
)

func SetupRequisitionApprovalsRoutes(r *gin.RouterGroup, handler *Handler) {
	api := r.Group("/requisitionapprovals")
	{
		api.GET("", handler.List)
		api.DELETE("/:id", handler.Approve)
		api.DELETE("/bulk-update", handler.ApproveBulk)
	}
}

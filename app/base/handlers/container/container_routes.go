package container

import (
	"github.com/gin-gonic/gin"
)

func SetupContainerRoutes(r *gin.RouterGroup, containerHandler *Handler) {
	api := r.Group("/containers")
	{
		api.DELETE("/bulk-update", containerHandler.DeleteBulk)
		api.GET("/getcontents", containerHandler.GetContentsByCode)
		api.POST("/marked-full", containerHandler.MarkedFull)

		api.GET("/reports/stock-levels", containerHandler.ReportStockLevels)
		api.GET("/reports/approvals", containerHandler.ReportApprovals)

		api.GET("", containerHandler.List)
		api.POST("", containerHandler.Create)
		api.GET("/:id", containerHandler.Get)
		api.PUT("/:id", containerHandler.Update)
		api.DELETE("/:id", containerHandler.Delete)
	}
}

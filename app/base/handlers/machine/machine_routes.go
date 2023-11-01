package machine

import (
	"github.com/gin-gonic/gin"
)

func SetupMachineRoutes(r *gin.RouterGroup, machineHandler *Handler) {
	api := r.Group("/machines")
	{
		api.DELETE("/bulk-update", machineHandler.DeleteBulk)
		api.GET("", machineHandler.List)
		api.POST("", machineHandler.Create)
		api.GET("/:id", machineHandler.Get)
		api.PUT("/:id", machineHandler.Update)
		api.DELETE("/:id", machineHandler.Delete)
	}
}

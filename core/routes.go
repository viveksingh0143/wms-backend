package core

import (
	"github.com/gin-gonic/gin"
	"star-wms/app/admin/handlers/permission"
	"star-wms/app/admin/handlers/plant"
	"star-wms/app/admin/handlers/role"
	"star-wms/app/admin/handlers/user"
	"star-wms/app/auth/handlers"
)

func SetupRoutes(r *gin.Engine, receiver *AppContainer) {
	api := r.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			v1.GET("/ping", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"message": "pong",
				})
			})

			authRoutes := api.Group("/auth")
			{
				handlers.SetupAuthRoutes(authRoutes, receiver.AuthHandler)
			}

			adminRoutes := api.Group("/admin")
			{
				permission.SetupPermissionRoutes(adminRoutes, receiver.PermissionHandler)
				role.SetupRoleRoutes(adminRoutes, receiver.RoleHandler)
				plant.SetupPlantRoutes(adminRoutes, receiver.PlantHandler)
				user.SetupUserRoutes(adminRoutes, receiver.UserHandler)
			}
		}
	}
}

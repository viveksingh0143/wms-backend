package core

import (
	"github.com/gin-gonic/gin"
	"star-wms/app/admin/handlers/permission"
	"star-wms/app/admin/handlers/plant"
	"star-wms/app/admin/handlers/role"
	"star-wms/app/admin/handlers/user"
	"star-wms/app/auth/handlers"
	"star-wms/app/base/handlers/category"
	"star-wms/app/base/handlers/container"
	"star-wms/app/base/handlers/customer"
	"star-wms/app/base/handlers/joborder"
	"star-wms/app/base/handlers/machine"
	"star-wms/app/base/handlers/product"
	"star-wms/app/base/handlers/store"
	"star-wms/core/middlewares"
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

			adminRoutes := api.Group("/admin", middlewares.AuthRequiredMiddleware(receiver.AuthService))
			{
				permission.SetupPermissionRoutes(adminRoutes, receiver.PermissionHandler)
				role.SetupRoleRoutes(adminRoutes, receiver.RoleHandler)
				plant.SetupPlantRoutes(adminRoutes, receiver.PlantHandler)
				user.SetupUserRoutes(adminRoutes, receiver.UserHandler)
			}

			baseRoutes := api.Group("/base", middlewares.AuthRequiredMiddleware(receiver.AuthService))
			{
				category.SetupCategoryRoutes(baseRoutes, receiver.CategoryHandler)
				plantBasedRoutes := baseRoutes.Group("/", middlewares.PlantRequiredMiddleware(receiver.PlantService))
				{
					product.SetupProductRoutes(plantBasedRoutes, receiver.ProductHandler)
					store.SetupStoreRoutes(plantBasedRoutes, receiver.StoreHandler)
					container.SetupContainerRoutes(plantBasedRoutes, receiver.ContainerHandler)
					machine.SetupMachineRoutes(plantBasedRoutes, receiver.MachineHandler)
					customer.SetupCustomerRoutes(plantBasedRoutes, receiver.CustomerHandler)
					joborder.SetupJobOrderRoutes(plantBasedRoutes, receiver.JobOrderHandler)
				}
			}
		}
	}
}

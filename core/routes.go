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
	"star-wms/app/base/handlers/outwardrequest"
	"star-wms/app/base/handlers/product"
	"star-wms/app/base/handlers/requisition"
	"star-wms/app/base/handlers/store"
	"star-wms/app/warehouse/handlers/batchlabel"
	"star-wms/app/warehouse/handlers/inventory"
	"star-wms/app/warehouse/handlers/requisitionapproval"
	"star-wms/app/warehouse/handlers/stockapproval"
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
		}

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
			product.SetupProductRoutes(baseRoutes, receiver.ProductHandler)

			plantBasedRoutes := baseRoutes.Group("/", middlewares.PlantRequiredMiddleware(receiver.PlantService))
			{
				store.SetupStoreRoutes(plantBasedRoutes, receiver.StoreHandler, receiver.StorelocationHandler)
				container.SetupContainerRoutes(plantBasedRoutes, receiver.ContainerHandler)
				machine.SetupMachineRoutes(plantBasedRoutes, receiver.MachineHandler)
				customer.SetupCustomerRoutes(plantBasedRoutes, receiver.CustomerHandler)
				joborder.SetupJoborderRoutes(plantBasedRoutes, receiver.JoborderHandler)
				requisition.SetupRequisitionRoutes(plantBasedRoutes, receiver.RequisitionHandler)
				outwardrequest.SetupOutwardrequestRoutes(plantBasedRoutes, receiver.OutwardrequestHandler)
			}
		}

		warehousePlantBasedRoutes := api.Group("/warehouse", middlewares.AuthRequiredMiddleware(receiver.AuthService), middlewares.PlantRequiredMiddleware(receiver.PlantService))
		{
			batchlabel.SetupBatchlabelRoutes(warehousePlantBasedRoutes, receiver.BatchlabelHandler, receiver.StickerHandler)
			inventory.SetupInventoryRoutes(warehousePlantBasedRoutes, receiver.InventoryHandler)
			stockapproval.SetupStockapprovalsRoutes(warehousePlantBasedRoutes, receiver.StockapprovalHandler)
			requisitionapproval.SetupRequisitionApprovalsRoutes(warehousePlantBasedRoutes, receiver.RequisitionApprovalHandler)
		}
	}
}

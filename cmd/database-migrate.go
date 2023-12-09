package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
	adminModels "star-wms/app/admin/models"
	baseModels "star-wms/app/base/models"
	warehouseModels "star-wms/app/warehouse/models"
	"star-wms/configs"
	"star-wms/core/migrators"
)

var databaseMigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "migrate database",
	Run: func(cmd *cobra.Command, args []string) {
		//if configs.AppCfg.TimeZone != "" {
		//	location, err := time.LoadLocation(configs.AppCfg.TimeZone)
		//	if err != nil {
		//		log.Fatal().Msgf("Failed to set time zone: %s", err)
		//		return
		//	}
		//	time.Local = location
		//}
		db, err := gorm.Open(configs.DBCfg.GetDatabaseConnection(), &gorm.Config{})
		if err != nil {
			log.Fatal().Msgf("Failed to connect to databaseMigrate: %v", err)
			return
		}

		migrator := migrators.NewDatabaseMigrator(db.Dialector.Name())
		if migrator == nil {
			log.Fatal().Msgf("unsupported database type: %s", db.Dialector.Name())
			return
		}
		err = migrator.Migrate(db, &adminModels.Permission{})
		if err != nil {
			log.Fatal().Msg("Could not migrate table for model: Permission")
			return
		}

		err = migrator.Migrate(db, &adminModels.Ability{})
		if err != nil {
			log.Fatal().Msg("Could not migrate table for model: Ability")
			return
		}

		err = migrator.Migrate(db, &adminModels.Role{})
		if err != nil {
			log.Fatal().Msg("Could not migrate table for model: Role")
			return
		}

		err = migrator.Migrate(db, &adminModels.Plant{})
		if err != nil {
			log.Fatal().Msg("Could not migrate table for model: Plant")
			return
		}

		err = migrator.Migrate(db, &adminModels.User{})
		if err != nil {
			log.Fatal().Msg("Could not migrate table for model: User")
			return
		}

		err = migrator.Migrate(db, &baseModels.Category{})
		if err != nil {
			log.Fatal().Msg("Could not migrate table for model: Category")
			return
		}

		err = migrator.Migrate(db, &baseModels.Product{})
		if err != nil {
			log.Fatal().Msg("Could not migrate table for model: Product")
			return
		}

		err = migrator.Migrate(db, &baseModels.ProductIngredient{})
		if err != nil {
			log.Fatal().Msg("Could not migrate table for model: ProductIngredient")
			return
		}

		err = migrator.Migrate(db, &baseModels.Store{})
		if err != nil {
			log.Fatal().Msg("Could not migrate table for model: Store")
			return
		}

		err = migrator.Migrate(db, &baseModels.Container{})
		if err != nil {
			log.Fatal().Msg("Could not migrate table for model: Container")
			return
		}

		err = migrator.Migrate(db, &baseModels.ContainerContent{})
		if err != nil {
			log.Fatal().Msg("Could not migrate table for model: ContainerContent")
			return
		}

		err = migrator.Migrate(db, &baseModels.Storelocation{})
		if err != nil {
			log.Fatal().Msg("Could not migrate table for model: Storelocation")
			return
		}

		err = migrator.Migrate(db, &baseModels.Machine{})
		if err != nil {
			log.Fatal().Msg("Could not migrate table for model: Machine")
			return
		}

		err = migrator.Migrate(db, &baseModels.Customer{})
		if err != nil {
			log.Fatal().Msg("Could not migrate table for model: Customer")
			return
		}

		err = migrator.Migrate(db, &baseModels.Joborder{})
		if err != nil {
			log.Fatal().Msg("Could not migrate table for model: Joborder")
			return
		}

		err = migrator.Migrate(db, &baseModels.JoborderItem{})
		if err != nil {
			log.Fatal().Msg("Could not migrate table for model: JoborderItem")
			return
		}

		err = migrator.Migrate(db, &baseModels.Requisition{})
		if err != nil {
			log.Fatal().Msg("Could not migrate table for model: Requisition")
			return
		}

		err = migrator.Migrate(db, &baseModels.RequisitionItem{})
		if err != nil {
			log.Fatal().Msg("Could not migrate table for model: RequisitionItem")
			return
		}

		err = migrator.Migrate(db, &baseModels.Outwardrequest{})
		if err != nil {
			log.Fatal().Msg("Could not migrate table for model: Outwardrequest")
			return
		}

		err = migrator.Migrate(db, &baseModels.OutwardrequestItem{})
		if err != nil {
			log.Fatal().Msg("Could not migrate table for model: OutwardrequestItem")
			return
		}

		err = migrator.Migrate(db, &warehouseModels.Batchlabel{})
		if err != nil {
			log.Fatal().Msg("Could not migrate table for model: Batchlabel")
			return
		}

		err = migrator.Migrate(db, &warehouseModels.Sticker{})
		if err != nil {
			log.Fatal().Msg("Could not migrate table for model: Sticker")
			return
		}

		err = migrator.Migrate(db, &warehouseModels.StickerItem{})
		if err != nil {
			log.Fatal().Msg("Could not migrate table for model: StickerItem")
			return
		}

		err = migrator.Migrate(db, &warehouseModels.Inventory{})
		if err != nil {
			log.Fatal().Msg("Could not migrate table for model: Inventory")
			return
		}

		err = migrator.Migrate(db, &warehouseModels.StockMovements{})
		if err != nil {
			log.Fatal().Msg("Could not migrate table for model: StockMovements")
			return
		}

		err = migrator.Migrate(db, &warehouseModels.RMBatch{})
		if err != nil {
			log.Fatal().Msg("Could not migrate table for model: RMBatch")
			return
		}

		err = migrator.Migrate(db, &warehouseModels.RMBatchTransaction{})
		if err != nil {
			log.Fatal().Msg("Could not migrate table for model: RMBatchTransaction")
			return
		}

		//err = db.AutoMigrate(
		//	&adminModels.Permission{},
		//	&adminModels.Ability{},
		//	&adminModels.Role{},
		//	&adminModels.Plant{},
		//	&adminModels.User{},
		//	&baseModels.Category{},
		//	&baseModels.Product{},
		//	&baseModels.ProductIngredient{},
		//	&baseModels.Store{},
		//	&baseModels.Container{},
		//	&baseModels.ContainerContent{},
		//	&baseModels.Storelocation{},
		//	&baseModels.Machine{},
		//	&baseModels.Customer{},
		//	&baseModels.Joborder{},
		//	&baseModels.JoborderItem{},
		//	&baseModels.Requisition{},
		//	&baseModels.RequisitionItem{},
		//	&baseModels.Outwardrequest{},
		//	&baseModels.OutwardrequestItem{},
		//	&warehouseModels.Batchlabel{},
		//	&warehouseModels.Sticker{},
		//	&warehouseModels.StickerItem{},
		//	&warehouseModels.Inventory{},
		//	&warehouseModels.StockMovements{},
		//	&warehouseModels.RMBatch{},
		//	&warehouseModels.RMBatchTransaction{},
		//)
		if err != nil {
			log.Fatal().Msgf("Database migration failed: %v", err)
			return
		} else {
			log.Info().Msg("Database migration successful.")
		}
	},
}

func init() {
	databaseCmd.AddCommand(databaseMigrateCmd)
}

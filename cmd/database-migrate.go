package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
	adminModels "star-wms/app/admin/models"
	baseModels "star-wms/app/base/models"
	"star-wms/configs"
	"time"
)

var databaseMigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "migrate database",
	Run: func(cmd *cobra.Command, args []string) {
		location, err := time.LoadLocation(configs.AppCfg.TimeZone)
		if err != nil {
			log.Fatal().Msgf("Failed to set time zone: %s", err)
			return
		}
		time.Local = location
		db, err := gorm.Open(configs.DBCfg.GetDatabaseConnection(), &gorm.Config{})
		if err != nil {
			log.Fatal().Msgf("Failed to connect to databaseMigrate: %v", err)
			return
		}

		err = db.AutoMigrate(
			&adminModels.Permission{}, &adminModels.Ability{}, &adminModels.Role{}, &adminModels.Plant{}, &adminModels.User{}, baseModels.Category{},
		)
		log.Info().Msg("Database auto migration for permissions")

		if err != nil {
			log.Fatal().Msgf("Could not migrate database: %v", err)
			return
		}
		log.Info().Msg("Database migration successful.")
	},
}

func init() {
	databaseCmd.AddCommand(databaseMigrateCmd)
}

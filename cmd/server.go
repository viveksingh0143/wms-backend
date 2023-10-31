package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
	"star-wms/configs"
	"star-wms/core"
	"star-wms/plugins/cache"
	"time"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the web application server",
	Run: func(cmd *cobra.Command, args []string) {
		location, err := time.LoadLocation(configs.AppCfg.TimeZone)
		if err != nil {
			log.Fatal().Msgf("Failed to set time zone: %s", err)
			return
		}
		time.Local = location

		customLogger := configs.ZeroLogGormLogger{Log: &log.Logger}
		db, err := gorm.Open(configs.DBCfg.GetDatabaseConnection(), &gorm.Config{
			Logger: customLogger,
		})
		if err != nil {
			log.Fatal().Msgf("Failed to connect to database: %v", err)
		}

		cacheManager, err := cache.NewCacheManager()
		if err != nil {
			log.Fatal().Msgf("Failed to create cache manager: %v", err)
		}

		appContainer := core.NewAppContainer(db, cacheManager)
		appContainer.RunServer()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}

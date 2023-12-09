package cmd

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
	"os"
	"star-wms/configs"
	"star-wms/core"
	"star-wms/plugins/cache"
	"strings"
)

var csvFilePath string
var resourceType string

var dataImportCmd = &cobra.Command{
	Use:   "data-import",
	Short: "Import data to the web application server",
	Run: func(cmd *cobra.Command, args []string) {
		// Set the time location
		//if configs.AppCfg.TimeZone != "" {
		//	location, err := time.LoadLocation(configs.AppCfg.TimeZone)
		//	if err != nil {
		//		log.Fatal().Msgf("Failed to set time zone: %s", err)
		//		return
		//	}
		//	time.Local = location
		//}

		// Database connection
		customLogger := configs.ZeroLogGormLogger{Log: &log.Logger}
		db, err := gorm.Open(configs.DBCfg.GetDatabaseConnection(), &gorm.Config{
			Logger: customLogger,
		})
		if err != nil {
			log.Fatal().Msgf("Failed to connect to database: %v", err)
		}

		// Cache manager setup
		cacheManager, err := cache.NewCacheManager()
		if err != nil {
			log.Fatal().Msgf("Failed to create cache manager: %v", err)
		}

		// Create the application container
		appContainer := core.NewAppContainer(db, cacheManager)

		// Check if the file path is provided
		if csvFilePath == "" {
			log.Fatal().Msg("CSV file path is required")
			return
		}

		resourceTypesAllowed := "-MATERIAL-PRODUCT-MACHINE-CUSTOMER-PERMISSION-ROLE-USER-CATEGORY-"

		if resourceType == "" || !strings.Contains(resourceTypesAllowed, fmt.Sprintf("-%s-", resourceType)) {
			log.Fatal().Msg("Resource info is required and it can be only MATERIAL, PRODUCT")
			return
		}

		// Check if the file exists
		if _, err := os.Stat(csvFilePath); os.IsNotExist(err) {
			log.Fatal().Msgf("CSV file does not exist: %s", csvFilePath)
			return
		}

		// Import data from CSV
		if resourceType == "MATERIAL" {
			_, err = appContainer.BulkService.ImportMaterialDataFromCSV(csvFilePath)
			if err != nil {
				log.Fatal().Msgf("Failed to import data from CSV: %v", err)
				return
			}
		} else if resourceType == "PRODUCT" {
			_, err = appContainer.BulkService.ImportProductDataFromCSV(csvFilePath)
			if err != nil {
				log.Fatal().Msgf("Failed to import data from CSV: %v", err)
				return
			}
		} else if resourceType == "MACHINE" {
			_, err = appContainer.BulkService.ImportProductDataFromCSV(csvFilePath)
			if err != nil {
				log.Fatal().Msgf("Failed to import data from CSV: %v", err)
				return
			}
		} else if resourceType == "CUSTOMER" {
			_, err = appContainer.BulkService.ImportProductDataFromCSV(csvFilePath)
			if err != nil {
				log.Fatal().Msgf("Failed to import data from CSV: %v", err)
				return
			}
		} else if resourceType == "PERMISSION" {
			_, err = appContainer.BulkService.ImportProductDataFromCSV(csvFilePath)
			if err != nil {
				log.Fatal().Msgf("Failed to import data from CSV: %v", err)
				return
			}
		} else if resourceType == "ROLE" {
			_, err = appContainer.BulkService.ImportProductDataFromCSV(csvFilePath)
			if err != nil {
				log.Fatal().Msgf("Failed to import data from CSV: %v", err)
				return
			}
		} else if resourceType == "USER" {
			_, err = appContainer.BulkService.ImportProductDataFromCSV(csvFilePath)
			if err != nil {
				log.Fatal().Msgf("Failed to import data from CSV: %v", err)
				return
			}
		} else if resourceType == "CATEGORY" {
			_, err = appContainer.BulkService.ImportProductDataFromCSV(csvFilePath)
			if err != nil {
				log.Fatal().Msgf("Failed to import data from CSV: %v", err)
				return
			}
		}
		log.Info().Msg("Data imported successfully from CSV")
	},
}

func init() {
	dataImportCmd.Flags().StringVarP(&csvFilePath, "file", "f", "", "Path to the CSV file to import")
	dataImportCmd.Flags().StringVarP(&resourceType, "resource", "r", "", "Type of resource")
	rootCmd.AddCommand(dataImportCmd)
}

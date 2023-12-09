package configs

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"sync"
)

var (
	once sync.Once
)

func InitRootConfig(cfgFile string) {
	once.Do(func() {
		if cfgFile != "" {
			viper.SetConfigFile(cfgFile)
		} else {
			home, err := os.UserHomeDir()
			cobra.CheckErr(err)
			viper.SetConfigType("yaml")
			viper.SetConfigName("star-wms.yaml")

			viper.AddConfigPath(home)
			viper.AddConfigPath(".")
			viper.AddConfigPath("/etc/star-wms/")
			viper.AddConfigPath("$HOME/.star-wms")
		}
		viper.AutomaticEnv() // read in environment variables that match
		if err := viper.ReadInConfig(); err == nil {
			log.Println("Using config file:", viper.ConfigFileUsed())
		} else {
			log.Fatalf("Error reading config file, %s", err)
		}

		// Initialize the sub-configurations
		InitAppConfig()
		InitLogConfig()
		InitDBConfig()
		InitServerConfig()
		InitAuthConfig()
	})
}

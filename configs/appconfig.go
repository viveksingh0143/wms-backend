package configs

import "github.com/spf13/viper"

type AppConfig struct {
	TimeZone string
	Debug    bool
}

var AppCfg *AppConfig

func InitAppConfig() {
	AppCfg = &AppConfig{
		TimeZone: viper.GetString("application.timezone"),
		Debug:    viper.GetBool("application.debug"),
	}
}

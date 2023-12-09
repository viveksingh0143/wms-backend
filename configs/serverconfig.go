package configs

import "github.com/spf13/viper"

type ServerConfig struct {
	Address string
	Port    int
}

var ServerCfg *ServerConfig

func InitServerConfig() {
	ServerCfg = &ServerConfig{
		Address: viper.GetString("rest-server.address"),
		Port:    viper.GetInt("rest-server.port"),
	}
}

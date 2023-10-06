package configs

import "github.com/spf13/viper"

type AuthConfig struct {
	ExpiryDuration     int64
	ExpiryLongDuration int64
	SecretKey          string
}

var AuthCfg *AuthConfig

func InitAuthConfig() {
	AuthCfg = &AuthConfig{
		ExpiryDuration:     viper.GetInt64("admin.expiry-duration"),
		ExpiryLongDuration: viper.GetInt64("admin.expiry-long-duration"),
		SecretKey:          viper.GetString("admin.secret-key"),
	}
}

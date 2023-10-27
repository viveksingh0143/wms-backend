package configs

import "github.com/spf13/viper"

type AuthConfig struct {
	ExpiryDuration            int64
	ExpiryLongDuration        int64
	RefreshExpiryDuration     int64
	RefreshExpiryLongDuration int64
	SecretKey                 string
}

var AuthCfg *AuthConfig

func InitAuthConfig() {
	AuthCfg = &AuthConfig{
		ExpiryDuration:            viper.GetInt64("auth.expiry-duration"),
		ExpiryLongDuration:        viper.GetInt64("auth.expiry-long-duration"),
		RefreshExpiryDuration:     viper.GetInt64("auth.refresh-expiry-duration"),
		RefreshExpiryLongDuration: viper.GetInt64("auth.refresh-expiry-long-duration"),
		SecretKey:                 viper.GetString("auth.secret-key"),
	}
}

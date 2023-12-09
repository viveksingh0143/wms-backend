package configs

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"net/url"
)

type DBConfig struct {
	Driver   string
	Host     string
	Port     int
	Username string
	Password string
	DBName   string
	TimeZone string
	LogQuery bool
}

var DBCfg *DBConfig

func InitDBConfig() {
	DBCfg = &DBConfig{
		Driver:   viper.GetString("database.driver"),
		Host:     viper.GetString("database.host"),
		Port:     viper.GetInt("database.port"),
		Username: viper.GetString("database.username"),
		Password: viper.GetString("database.password"),
		DBName:   viper.GetString("database.dbname"),
		TimeZone: viper.GetString("database.timezone"),
		LogQuery: viper.GetBool("database.show-query"),
	}
}

func (c *DBConfig) GetDatabaseConnection() gorm.Dialector {
	dbPassword := url.QueryEscape(DBCfg.Password)
	log.Debug().Msgf("DB Port: %d", DBCfg.Port)
	switch c.Driver {
	case "mysql":
		timeZone := url.QueryEscape(DBCfg.TimeZone)
		dns := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=%s", DBCfg.Username, dbPassword, DBCfg.Host, DBCfg.Port, DBCfg.DBName, timeZone)
		return mysql.Open(dns)
	case "mssql":
		timeZone := url.QueryEscape(DBCfg.TimeZone)
		dns := fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s&charset=utf8mb4&parseTime=True&loc=%s", DBCfg.Username, dbPassword, DBCfg.Host, DBCfg.Port, DBCfg.DBName, timeZone)
		return sqlserver.Open(dns)
	case "postgres":
		dns := fmt.Sprintf("user=%s password=%s host=%s dbname=%s port=%d sslmode=disable TimeZone=%s", DBCfg.Username, dbPassword, DBCfg.Host, DBCfg.DBName, DBCfg.Port, DBCfg.TimeZone)
		return postgres.Open(dns)
	case "sqlite3":
		dns := c.DBName
		return sqlite.Open(dns)
	default:
		log.Fatal().Msgf("Unsupported database driver: %s", c.Driver)
		return nil
	}
}

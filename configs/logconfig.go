package configs

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
)

type LogConfig struct {
	ConsoleLog        bool
	FileLog           bool
	Level             string
	FilePath          string
	FileMaxSizeInMB   uint
	MaxBackupCount    uint
	MaxAgeInDays      uint
	CompressBackupLog bool
}

var LogCfg *LogConfig

func InitLogConfig() {
	LogCfg = &LogConfig{
		ConsoleLog:        viper.GetBool("logger.console-log"),
		FileLog:           viper.GetBool("logger.file-log"),
		Level:             viper.GetString("logger.level"),
		FilePath:          viper.GetString("logger.file-path"),
		FileMaxSizeInMB:   viper.GetUint("logger.max-size"),
		MaxBackupCount:    viper.GetUint("logger.max-backups"),
		MaxAgeInDays:      viper.GetUint("logger.max-age"),
		CompressBackupLog: viper.GetBool("logger.compress"),
	}

	LogCfg.InitializeLogger()
}

func (lc *LogConfig) InitializeLogger() {
	var writers []io.Writer

	if lc.ConsoleLog {
		consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}
		writers = append(writers, consoleWriter)
	}

	if lc.FileLog {
		fileWriter := &lumberjack.Logger{
			Filename:   lc.FilePath,
			MaxSize:    int(lc.FileMaxSizeInMB),
			MaxBackups: int(lc.MaxBackupCount),
			MaxAge:     int(lc.MaxAgeInDays),
			Compress:   lc.CompressBackupLog,
		}
		writers = append(writers, fileWriter)
	}

	multi := zerolog.MultiLevelWriter(writers...)
	log.Logger = zerolog.New(multi).With().Timestamp().Logger()

	// Set log level
	level, err := zerolog.ParseLevel(lc.Level)
	if err != nil {
		log.Error().Err(err).Msg("Invalid log level; using 'info' as default")
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	// Example log
	log.Info().Msg("logger is configured.")
}

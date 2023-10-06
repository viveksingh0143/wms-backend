package configs

import (
	"context"
	"github.com/rs/zerolog"
	"gorm.io/gorm/logger"
	"time"
)

type ZeroLogGormLogger struct {
	Log *zerolog.Logger
}

func (l ZeroLogGormLogger) LogMode(_ logger.LogLevel) logger.Interface {
	return l
}

func (l ZeroLogGormLogger) Info(_ context.Context, msg string, data ...interface{}) {
	l.Log.Info().Fields(parseData(data)).Msg(msg)
}

func (l ZeroLogGormLogger) Warn(_ context.Context, msg string, data ...interface{}) {
	l.Log.Warn().Fields(parseData(data)).Msg(msg)
}

func (l ZeroLogGormLogger) Error(_ context.Context, msg string, data ...interface{}) {
	l.Log.Error().Fields(parseData(data)).Msg(msg)
}

func (l ZeroLogGormLogger) Trace(_ context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	message, rows := fc()
	if err != nil {
		l.Log.Error().Err(err).Fields(map[string]interface{}{
			"elapsed": elapsed,
			"rows":    rows,
			"sql":     message,
		}).Msg("gorm query")
	} else {
		l.Log.Debug().Fields(map[string]interface{}{
			"elapsed": elapsed,
			"rows":    rows,
			"sql":     message,
		}).Msg("gorm query")
	}
}

// Helper function to convert variadic interface{} to a map for ZeroLog Fields
func parseData(data []interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for i := 0; i < len(data); i += 2 {
		key, ok := data[i].(string)
		if !ok {
			continue
		}
		value := data[i+1]
		result[key] = value
	}
	return result
}

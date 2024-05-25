package env

import (
	"os"

	"github.com/kirillmc/auth/internal/config"
)

const (
	logLevelEnvName = "LOG_LEVEL"
)

type loggerConfig struct {
	logLevel string
}

func NewLoggerConfig() config.LoggerConfig {
	logLevel := os.Getenv(logLevelEnvName)
	if len(logLevel) == 0 {
		return &loggerConfig{logLevel: "warn"}
	}

	return &loggerConfig{logLevel: logLevel}
}

func (cfg *loggerConfig) LogLevel() string {
	return cfg.logLevel
}

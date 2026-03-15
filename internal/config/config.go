package config

import (
	"fmt"
	"strings"

	"github.com/caarlos0/env/v11"
)

const (
	LogLevelDebug string = "debug"
	LogLevelInfo         = "info"
	LogLevelWarn         = "warn"
	LogLevelError        = "error"
	LogLevelPanic        = "panic"
	LogLevelFatal        = "fatal"
)

const (
	DefaultBaseURL       = "http://localhost:8080"
	DefaultServerAddress = ":8080"
	DefaultEnableLogs    = false
	DefaultLogLevel      = LogLevelInfo
)

const (
	ServerAddressEnvName = "SERVER_ADDRESS"
	BaseURLEnvName       = "BASE_URL"
	EnableLogsEnvName    = "ENABLE_LOGS"
	LogLevelEnvName      = "LOG_LEVEL"
	LogFileEnvName       = "LOG_FILE"
)

type Config struct {
	ServerAddress string  `env:"SERVER_ADDRESS" default:":8080"`
	BaseURL       string  `env:"BASE_URL" default:"http://localhost:8080"`
	EnableLogs    bool    `env:"ENABLE_LOGS" default:"false"`
	LogLevel      string  `env:"LOG_LEVEL" default:"info"`
	LogFile       *string `env:"LOG_FILE" default:""`
}

func ParseEnv(prefix string, cfg *Config) error {
	if cfg == nil {
		return nil
	}
	err := env.ParseWithOptions(cfg, env.Options{
		Prefix:              prefix,
		TagName:             "env",
		DefaultValueTagName: "default",
	})
	if err != nil {
		return err
	}
	return nil
}

func getEnvNameWithPrefix(prefix, envName string) string {
	if prefix != "" {
		return strings.ToUpper(fmt.Sprintf("%s_%s", prefix, envName))
	}
	return strings.ToUpper(envName)
}

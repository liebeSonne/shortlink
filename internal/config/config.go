package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/caarlos0/env/v11"
)

const (
	LogLevelDebug = "debug"
	LogLevelInfo  = "info"
	LogLevelWarn  = "warn"
	LogLevelError = "error"
	LogLevelPanic = "panic"
	LogLevelFatal = "fatal"
)

const (
	DefaultBaseURL            = "http://localhost:8080"
	DefaultServerAddress      = ":8080"
	DefaultEnableLogs         = false
	DefaultLogLevel           = LogLevelInfo
	DefaultAuthCookieTokenKey = "session_token"
	DefaultAuthSecretKey      = "secret-key-123"
	DefaultAuthTokenExpires   = time.Hour * 24
)

const (
	ServerAddressEnvName      = "SERVER_ADDRESS"
	BaseURLEnvName            = "BASE_URL"
	EnableLogsEnvName         = "ENABLE_LOGS"
	LogLevelEnvName           = "LOG_LEVEL"
	LogFileEnvName            = "LOG_FILE"
	FileStoragePathEnvName    = "FILE_STORAGE_PATH"
	DatabaseDSNEnvName        = "DATABASE_DSN"
	AuthCookieTokenKeyEnvName = "AUTH_COOKIE_TOKEN_KEY"
	AuthSecretKeyEnvName      = "AUTH_SECRET_KEY"
	AuthTokenExpiresEnvName   = "AUTH_TOKEN_EXPIRE"
)

type Config struct {
	ServerAddress   string  `env:"SERVER_ADDRESS" default:":8080"`
	BaseURL         string  `env:"BASE_URL" default:"http://localhost:8080"`
	EnableLogs      bool    `env:"ENABLE_LOGS" default:"false"`
	LogLevel        string  `env:"LOG_LEVEL" default:"info"`
	LogFile         *string `env:"LOG_FILE" default:""`
	FileStoragePath *string `env:"FILE_STORAGE_PATH" default:""`
	DatabaseDSN     *string `env:"DATABASE_DSN" default:""`

	AuthCookieTokenKey string        `env:"AUTH_COOKIE_TOKEN_KEY" default:"session_token"`
	AuthSecretKey      string        `env:"AUTH_SECRET_KEY" default:"secret-key-123"`
	AuthTokenExpires   time.Duration `env:"AUTH_TOKEN_EXPIRE" default:"24h"`
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

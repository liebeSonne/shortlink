package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/caarlos0/env/v11"
)

// Значения уровней логирования.
const (
	LogLevelDebug = "debug"
	LogLevelInfo  = "info"
	LogLevelWarn  = "warn"
	LogLevelError = "error"
	LogLevelPanic = "panic"
	LogLevelFatal = "fatal"
)

// Значения по умолчанию.
const (
	DefaultBaseURL            = "http://localhost:8080"
	DefaultServerAddress      = ":8080"
	DefaultEnableLogs         = false
	DefaultLogLevel           = LogLevelInfo
	DefaultAuthCookieTokenKey = "session_token"
	DefaultAuthSecretKey      = "secret-key-123"
	DefaultAuthTokenExpires   = time.Hour * 24
)

// Названия настроек в переменных окружения.
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
	AuditFileEnvName          = "AUDIT_FILE"
	AuditURLEnvName           = "AUDIT_URL"
)

// Config - настройки.
type Config struct {
	ServerAddress   string  `env:"SERVER_ADDRESS" default:":8080"`           // Адрес сервера.
	BaseURL         string  `env:"BASE_URL" default:"http://localhost:8080"` // Базовая ссылка создаваемых сокращенных ссылок.
	EnableLogs      bool    `env:"ENABLE_LOGS" default:"false"`              // Включение логирования.
	LogLevel        string  `env:"LOG_LEVEL" default:"info"`                 // Уровень логирования.
	LogFile         *string `env:"LOG_FILE" default:""`                      // Файл для сохранения логов.
	FileStoragePath *string `env:"FILE_STORAGE_PATH" default:""`             // Файловое хранилище для сокращенных ссылок.
	DatabaseDSN     *string `env:"DATABASE_DSN" default:""`                  // Параметры соединения с базой данных.

	AuthCookieTokenKey string        `env:"AUTH_COOKIE_TOKEN_KEY" default:"session_token"` // Название токена авторизации с cookie.
	AuthSecretKey      string        `env:"AUTH_SECRET_KEY" default:"secret-key-123"`      // Секретный код для подписи токена.
	AuthTokenExpires   time.Duration `env:"AUTH_TOKEN_EXPIRE" default:"24h"`               // Время жизни токена.

	AuditFile *string `env:"AUDIT_FILE" default:""` // Файл для сохранения данных аудита.
	AuditURL  *string `env:"AUDIT_URL" default:""`  // Ссылка для сохранения данных аудита.
}

// ParseEnv - парсинг настроек из переменных окружения.
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

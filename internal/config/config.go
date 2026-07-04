package config

import (
	"encoding/json"
	"fmt"
	"os"
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
	DefaultEnableHTTPS        = false
)

var defaultConfig = Config{
	BaseURL:            DefaultBaseURL,
	ServerAddress:      DefaultServerAddress,
	EnableLogs:         DefaultEnableLogs,
	LogLevel:           DefaultLogLevel,
	AuthCookieTokenKey: DefaultAuthCookieTokenKey,
	AuthSecretKey:      DefaultAuthSecretKey,
	AuthTokenExpires:   DefaultAuthTokenExpires,
	EnableHTTPS:        DefaultEnableHTTPS,
}

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
	EnableHTTPSEnvName        = "ENABLE_HTTPS"
	ConfigFileEnvName         = "CONFIG"
)

var allEnvNames = []string{
	ServerAddressEnvName,
	BaseURLEnvName,
	EnableLogsEnvName,
	LogLevelEnvName,
	LogFileEnvName,
	FileStoragePathEnvName,
	DatabaseDSNEnvName,
	AuthCookieTokenKeyEnvName,
	AuthSecretKeyEnvName,
	AuthTokenExpiresEnvName,
	AuditFileEnvName,
	AuditURLEnvName,
	EnableHTTPSEnvName,
	ConfigFileEnvName,
}

// Config - настройки.
type Config struct {
	ServerAddress   string  `env:"SERVER_ADDRESS" default:":8080" json:"server_address"`     // Адрес сервера.
	BaseURL         string  `env:"BASE_URL" default:"http://localhost:8080" json:"base_url"` // Базовая ссылка создаваемых сокращенных ссылок.
	EnableLogs      bool    `env:"ENABLE_LOGS" default:"false" json:"enable_logs"`           // Включение логирования.
	LogLevel        string  `env:"LOG_LEVEL" default:"info" json:"log_level"`                // Уровень логирования.
	LogFile         *string `env:"LOG_FILE" default:"" json:"log_file"`                      // Файл для сохранения логов.
	FileStoragePath *string `env:"FILE_STORAGE_PATH" default:"" json:"file_storage_path"`    // Файловое хранилище для сокращенных ссылок.
	DatabaseDSN     *string `env:"DATABASE_DSN" default:"" json:"database_dsn"`              // Параметры соединения с базой данных.

	AuthCookieTokenKey string        `env:"AUTH_COOKIE_TOKEN_KEY" default:"session_token" json:"auth_cookie_token_key"` // Название токена авторизации с cookie.
	AuthSecretKey      string        `env:"AUTH_SECRET_KEY" default:"secret-key-123" json:"auth_secret_key"`            // Секретный код для подписи токена.
	AuthTokenExpires   time.Duration `env:"AUTH_TOKEN_EXPIRE" default:"24h" json:"auth_token_expires"`                  // Время жизни токена.

	AuditFile *string `env:"AUDIT_FILE" default:"" json:"audit_file"` // Файл для сохранения данных аудита.
	AuditURL  *string `env:"AUDIT_URL" default:"" json:"audit_url"`   // Ссылка для сохранения данных аудита.

	EnableHTTPS bool `env:"ENABLE_HTTPS" default:"false" json:"enable_https"` // Включение HTTPS в веб-сервере

	ConfigFile *string `env:"CONFIG" default:""` // json-файл конфигурации
}

func MakeModConfig(c Config, f func(c *Config)) Config {
	f(&c)
	return c
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

func ParseFromJSON(path string, cfg *Config) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read config file %s: %w", path, err)
	}

	err = json.Unmarshal(data, cfg)
	if err != nil {
		return fmt.Errorf("failed to parse config file %s: %w", path, err)
	}

	return nil
}

func getEnvNameWithPrefix(prefix, envName string) string {
	if prefix != "" {
		return strings.ToUpper(fmt.Sprintf("%s_%s", prefix, envName))
	}
	return strings.ToUpper(envName)
}

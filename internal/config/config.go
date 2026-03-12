package config

import (
	"fmt"
	"strings"

	"github.com/caarlos0/env/v11"
)

const DefaultBaseURL = "http://localhost:8080"
const DefaultServerAddress = ":8080"
const DefaultEnableLogs = false

const ServerAddressEnvName = "SERVER_ADDRESS"
const BaseURLEnvName = "BASE_URL"
const EnableLogsEnvName = "ENABLE_LOGS"

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS" default:":8080"`
	BaseURL       string `env:"BASE_URL" default:"http://localhost:8080"`
	EnableLogs    bool   `env:"ENABLE_LOGS" default:"false"`
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

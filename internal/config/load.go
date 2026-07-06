package config

import (
	"fmt"
	"os"

	"github.com/caarlos0/env/v11"
)

// LoadConfig - загружает настройки из переменных окружения (первый приоритет) и из флагов (второй приоритет) и из файла (третий приоритет).
func LoadConfig(appID, envPrefix string) (Config, error) {
	// загружаем настройки из флагов
	fCfg := flagsConfig{}
	err := parseFlagsConfig(appID, &fCfg, true)
	if err != nil {
		return Config{}, fmt.Errorf("error parsing flags: %w", err)
	}

	prefix := getEnvNameWithPrefix(envPrefix, "")

	// определяем путь к файлу конфига
	var configFile string
	if fCfg.ConfigFile != nil && *fCfg.ConfigFile != "" {
		configFile = *fCfg.ConfigFile
	}
	if name, ok := os.LookupEnv(ConfigFileEnvName); ok {
		configFile = name
	}

	// загружаем настройки из default
	cfg := Config{}
	err = env.ParseWithOptions(&cfg, env.Options{
		Prefix:              prefix,
		TagName:             "not-use-env",
		DefaultValueTagName: "default",
	})
	if err != nil {
		return Config{}, fmt.Errorf("error parsing config: %w", err)
	}

	// загружаем настройки из JSON файла (не переопределенные настройки default останутся)
	if configFile != "" {
		err = ParseFromJSON(configFile, &cfg)
		if err != nil {
			return Config{}, fmt.Errorf("error parsing json-file config: %w", err)
		}
	}

	// загружаем настройки из flags
	mergeFlagsConfig(fCfg, &cfg, allEnvNames)

	// загружаем настройки из env (не переопределенные настройки из файла и flags останутся)
	err = env.ParseWithOptions(&cfg, env.Options{
		Prefix:              prefix,
		TagName:             "env",
		DefaultValueTagName: "not-use",
	})
	if err != nil {
		return Config{}, fmt.Errorf("error parsing env: %w", err)
	}

	return cfg, nil
}

// mergeFlagsConfig - слияние настроек с настройками из флага для перечисленных переменных окружения.
func mergeFlagsConfig(fCfg flagsConfig, cfg *Config, envNames []string) {
	if cfg == nil {
		return
	}

	for _, envName := range envNames {
		switch envName {
		case ServerAddressEnvName:
			if fCfg.ServerAddress != nil {
				cfg.ServerAddress = *fCfg.ServerAddress
			}
		case BaseURLEnvName:
			if fCfg.BaseURL != nil {
				cfg.BaseURL = *fCfg.BaseURL
			}
		case EnableLogsEnvName:
			if fCfg.EnableLogs != nil {
				cfg.EnableLogs = *fCfg.EnableLogs
			}
		case LogLevelEnvName:
			if fCfg.LogLevel != nil {
				cfg.LogLevel = *fCfg.LogLevel
			}
		case LogFileEnvName:
			if fCfg.LogFile != nil && *fCfg.LogFile != "" {
				cfg.LogFile = fCfg.LogFile
			}
		case FileStoragePathEnvName:
			if fCfg.FileStoragePath != nil && *fCfg.FileStoragePath != "" {
				cfg.FileStoragePath = fCfg.FileStoragePath
			}
		case DatabaseDSNEnvName:
			if fCfg.DatabaseDSN != nil && *fCfg.DatabaseDSN != "" {
				cfg.DatabaseDSN = fCfg.DatabaseDSN
			}
		case AuditFileEnvName:
			if fCfg.AuditFile != nil && *fCfg.AuditFile != "" {
				cfg.AuditFile = fCfg.AuditFile
			}
		case AuditURLEnvName:
			if fCfg.AuditURL != nil && *fCfg.AuditURL != "" {
				cfg.AuditURL = fCfg.AuditURL
			}
		case EnableHTTPSEnvName:
			if fCfg.EnableHTTPS != nil {
				cfg.EnableHTTPS = *fCfg.EnableHTTPS
			}
		case TLSCertFileEnvName:
			if fCfg.TLSCertFile != nil && *fCfg.TLSCertFile != "" {
				cfg.TLSCertFile = *fCfg.TLSCertFile
			}
		case TLSKeyFileEnvName:
			if fCfg.TLSKeyFile != nil && *fCfg.TLSKeyFile != "" {
				cfg.TLSKeyFile = *fCfg.TLSKeyFile
			}
		case ConfigFileEnvName:
			if fCfg.ConfigFile != nil && *fCfg.ConfigFile != "" {
				cfg.ConfigFile = fCfg.ConfigFile
			}
		}
	}
}

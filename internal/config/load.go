package config

import (
	"fmt"
	"maps"

	"github.com/caarlos0/env/v11"
)

// LoadConfig - Load config from env (first priority) and from flags (second priority)
func LoadConfig(appID, envPrefix string) (Config, error) {
	fCfg := flagsConfig{}
	err := parseFlagsConfig(appID, &fCfg, true)
	if err != nil {
		return Config{}, fmt.Errorf("error parsing flags: %w", err)
	}

	prefix := getEnvNameWithPrefix(envPrefix, "")

	tagNameToEnvName := map[string]string{
		prefix + ServerAddressEnvName:   ServerAddressEnvName,
		prefix + BaseURLEnvName:         BaseURLEnvName,
		prefix + EnableLogsEnvName:      EnableLogsEnvName,
		prefix + LogLevelEnvName:        LogLevelEnvName,
		prefix + LogFileEnvName:         LogFileEnvName,
		prefix + FileStoragePathEnvName: FileStoragePathEnvName,
		prefix + DatabaseDSNEnvName:     DatabaseDSNEnvName,
	}

	onSetHook := func(tag string, value interface{}, isDefault bool) {
		if !isDefault {
			delete(tagNameToEnvName, tag)
		}
	}

	cfg := Config{}
	err = env.ParseWithOptions(&cfg, env.Options{
		OnSet:               onSetHook,
		Prefix:              prefix,
		TagName:             "env",
		DefaultValueTagName: "default",
	})
	if err != nil {
		return Config{}, fmt.Errorf("error parsing env: %w", err)
	}

	envNames := make([]string, 0)
	for v := range maps.Values(tagNameToEnvName) {
		envNames = append(envNames, v)
	}

	mergeFlagsConfig(fCfg, &cfg, envNames)

	return cfg, nil
}

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
		}

	}
}

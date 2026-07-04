package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Названия флагов.
const (
	ServerAddressFlagName   = "a"
	BaseURLFlagName         = "b"
	EnableLogsFlagName      = "l"
	LogLevelFlagName        = "ll"
	LogFileFlagName         = "lf"
	FileStoragePathFlagName = "f"
	DatabaseDSNFlagName     = "d"
	AuditFileFlagName       = "audit-file"
	AuditURLFlagName        = "audit-url"
	EnableHTTPSFlagName     = "s"
	ConfigShortFlagName     = "c"
	ConfigLongFlagName      = "config"
)

// Ошибки разора флагов.
var (
	ErrInvalidFlagValue            = errors.New("invalid flag value")
	ErrInvalidDefaultServerAddress = errors.New("invalid default server address")
)

// flagsConfig - настройки получаемые из флагов.
type flagsConfig struct {
	ServerAddress   *string
	BaseURL         *string
	EnableLogs      *bool
	LogLevel        *string
	LogFile         *string
	FileStoragePath *string
	DatabaseDSN     *string
	AuditFile       *string
	AuditURL        *string
	EnableHTTPS     *bool
	ConfigFile      *string
}

func makeModFlagConfig(c flagsConfig, f func(c *flagsConfig)) flagsConfig {
	f(&c)
	return c
}

// parseFlags - получение настроек из флагов.
func parseFlags(appID string, config *Config) error {
	flagsConf := flagsConfig{}
	err := parseFlagsConfig(appID, &flagsConf, false)
	if err != nil {
		return err
	}

	if config != nil {
		if flagsConf.ServerAddress != nil {
			config.ServerAddress = *flagsConf.ServerAddress
		}
		if flagsConf.BaseURL != nil {
			config.BaseURL = *flagsConf.BaseURL
		}
		if flagsConf.EnableLogs != nil {
			config.EnableLogs = *flagsConf.EnableLogs
		}
		if flagsConf.LogLevel != nil {
			config.LogLevel = *flagsConf.LogLevel
		}
		if flagsConf.LogFile != nil {
			config.LogFile = flagsConf.LogFile
		}
		if flagsConf.FileStoragePath != nil {
			config.FileStoragePath = flagsConf.FileStoragePath
		}
		if flagsConf.DatabaseDSN != nil {
			config.DatabaseDSN = flagsConf.DatabaseDSN
		}
		if flagsConf.AuditFile != nil {
			config.AuditFile = flagsConf.AuditFile
		}
		if flagsConf.AuditURL != nil {
			config.AuditURL = flagsConf.AuditURL
		}
		if flagsConf.EnableHTTPS != nil {
			config.EnableHTTPS = *flagsConf.EnableHTTPS
		}
		if flagsConf.ConfigFile != nil {
			config.ConfigFile = flagsConf.ConfigFile
		}
	}

	return nil
}

// parseFlagsConfig - инициализация настроек флагов из флагов.
func parseFlagsConfig(appID string, config *flagsConfig, justIfSet bool) error {
	fs := flag.NewFlagSet(appID, flag.ContinueOnError)

	serverAddress := address{}
	err := serverAddress.Set(DefaultServerAddress)
	if err != nil {
		return errors.Join(err, ErrInvalidDefaultServerAddress)
	}

	fs.Var(&serverAddress, ServerAddressFlagName, "address and port to run server")
	baseURL := fs.String(BaseURLFlagName, DefaultBaseURL, "address and port for output short url")
	enableLogs := fs.Bool(EnableLogsFlagName, DefaultEnableLogs, "enable output logs")
	logLevel := fs.String(LogLevelFlagName, DefaultLogLevel, "log level")
	logFile := fs.String(LogFileFlagName, "", "log file")
	fileStoragePath := fs.String(FileStoragePathFlagName, "", "file storage path")
	databaseDSN := fs.String(DatabaseDSNFlagName, "", "database DSN")
	auditFile := fs.String(AuditFileFlagName, "", "audit file")
	auditURL := fs.String(AuditURLFlagName, "", "audit URL")
	enableHTTPS := fs.Bool(EnableHTTPSFlagName, DefaultEnableHTTPS, "enable HTTPS on web-server")
	var configFile string
	fs.StringVar(&configFile, ConfigShortFlagName, "", "config json file")
	fs.StringVar(&configFile, ConfigLongFlagName, "", "config json file")

	err = fs.Parse(os.Args[1:])
	if err != nil {
		return fmt.Errorf("error parsing config flags: %w", err)
	}

	if config == nil {
		return nil
	}

	if justIfSet {
		isSetFlagMap := map[string]bool{
			ServerAddressFlagName:   false,
			BaseURLFlagName:         false,
			EnableLogsFlagName:      false,
			LogLevelFlagName:        false,
			LogFileFlagName:         false,
			FileStoragePathFlagName: false,
			DatabaseDSNFlagName:     false,
			AuditFileFlagName:       false,
			AuditURLFlagName:        false,
			EnableHTTPSFlagName:     false,
			ConfigShortFlagName:     false,
			ConfigLongFlagName:      false,
		}

		fs.Visit(func(f *flag.Flag) {
			isSetFlagMap[f.Name] = true
		})

		if isSet, ok := isSetFlagMap[ServerAddressFlagName]; ok && isSet {
			addr := serverAddress.String()
			config.ServerAddress = &addr
		}
		if isSet, ok := isSetFlagMap[BaseURLFlagName]; ok && isSet {
			config.BaseURL = baseURL
		}
		if isSet, ok := isSetFlagMap[EnableLogsFlagName]; ok && isSet {
			config.EnableLogs = enableLogs
		}
		if isSet, ok := isSetFlagMap[LogLevelFlagName]; ok && isSet {
			config.LogLevel = logLevel
		}
		if isSet, ok := isSetFlagMap[LogFileFlagName]; ok && isSet {
			if logFile != nil && *logFile != "" {
				config.LogFile = logFile
			}
		}
		if isSet, ok := isSetFlagMap[FileStoragePathFlagName]; ok && isSet {
			if fileStoragePath != nil && *fileStoragePath != "" {
				config.FileStoragePath = fileStoragePath
			}
		}
		if isSet, ok := isSetFlagMap[DatabaseDSNFlagName]; ok && isSet {
			if databaseDSN != nil && *databaseDSN != "" {
				config.DatabaseDSN = databaseDSN
			}
		}
		if isSet, ok := isSetFlagMap[AuditFileFlagName]; ok && isSet {
			if auditFile != nil && *auditFile != "" {
				config.AuditFile = auditFile
			}
		}
		if isSet, ok := isSetFlagMap[AuditURLFlagName]; ok && isSet {
			if auditURL != nil && *auditURL != "" {
				config.AuditURL = auditURL
			}
		}
		if isSet, ok := isSetFlagMap[EnableHTTPSFlagName]; ok && isSet {
			config.EnableHTTPS = enableHTTPS
		}
		if isSet, ok := isSetFlagMap[ConfigShortFlagName]; ok && isSet {
			if configFile != "" {
				config.ConfigFile = &configFile
			}
		}
		if isSet, ok := isSetFlagMap[ConfigLongFlagName]; ok && isSet {
			if configFile != "" {
				config.ConfigFile = &configFile
			}
		}
	} else {
		addr := serverAddress.String()
		config.ServerAddress = &addr
		config.BaseURL = baseURL
		config.EnableLogs = enableLogs
		config.LogLevel = logLevel
		if logFile != nil && *logFile != "" {
			config.LogFile = logFile
		}
		if fileStoragePath != nil && *fileStoragePath != "" {
			config.FileStoragePath = fileStoragePath
		}
		if databaseDSN != nil && *databaseDSN != "" {
			config.DatabaseDSN = databaseDSN
		}
		if auditFile != nil && *auditFile != "" {
			config.AuditFile = auditFile
		}
		if auditURL != nil && *auditURL != "" {
			config.AuditURL = auditURL
		}
		config.EnableHTTPS = enableHTTPS
		if configFile != "" {
			config.ConfigFile = &configFile
		}
	}

	return nil
}

// Address - параметры настройки адреса.
type address struct {
	Host string
	Port int
}

func (a *address) String() string {
	return fmt.Sprintf("%s:%d", a.Host, a.Port)
}
func (a *address) Set(flagValue string) error {
	params := strings.Split(flagValue, ":")

	if len(params) != 2 {
		return ErrInvalidFlagValue
	}

	port, err := strconv.Atoi(params[1])
	if err != nil {
		return errors.Join(fmt.Errorf("error on atoi port: %w", err), ErrInvalidFlagValue)
	}

	a.Host = params[0]
	a.Port = port
	return nil
}

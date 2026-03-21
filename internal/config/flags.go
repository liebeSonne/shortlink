package config

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	ServerAddressFlagName   = "a"
	BaseURLFlagName         = "b"
	EnableLogsFlagName      = "l"
	LogLevelFlagName        = "ll"
	LogFileFlagName         = "lf"
	FileStoragePathFlagName = "f"
)

var ErrInvalidFlagValue = errors.New("invalid flag value")
var ErrInvalidDefaultServerAddress = errors.New("invalid default server address")

type flagsConfig struct {
	ServerAddress   *string
	BaseURL         *string
	EnableLogs      *bool
	LogLevel        *string
	LogFile         *string
	FileStoragePath *string
}

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
	}

	return nil
}

func parseFlagsConfig(appID string, config *flagsConfig, justIfSet bool) error {
	fs := flag.NewFlagSet(appID, flag.ContinueOnError)

	serverAddress := address{}
	err := serverAddress.Set(DefaultServerAddress)
	if err != nil {
		log.Printf("invalid default server address: %v", err)
		return ErrInvalidDefaultServerAddress
	}

	fs.Var(&serverAddress, ServerAddressFlagName, "address and port to run server")
	baseURL := fs.String(BaseURLFlagName, DefaultBaseURL, "address and port for output short url")
	enableLogs := fs.Bool(EnableLogsFlagName, DefaultEnableLogs, "enable output logs")
	logLevel := fs.String(LogLevelFlagName, DefaultLogLevel, "log level")
	logFile := fs.String(LogFileFlagName, "", "log file")
	fileStoragePath := fs.String(FileStoragePathFlagName, "", "file storage path")

	err = fs.Parse(os.Args[1:])
	if err != nil {
		log.Printf("error parsing config flags: %v", err)
		return err
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
	}

	return nil
}

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
		log.Printf("error on atoi port: %v\n", err)
		return ErrInvalidFlagValue
	}

	a.Host = params[0]
	a.Port = port
	return nil
}

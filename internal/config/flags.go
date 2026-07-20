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
	ServerAddressFlagName     = "a"
	GRPCServerAddressFlagName = "g"
	BaseURLFlagName           = "b"
	EnableLogsFlagName        = "l"
	LogLevelFlagName          = "ll"
	LogFileFlagName           = "lf"
	FileStoragePathFlagName   = "f"
	DatabaseDSNFlagName       = "d"
	AuditFileFlagName         = "audit-file"
	AuditURLFlagName          = "audit-url"
	EnableHTTPSFlagName       = "s"
	TLSCertFileFlagName       = "tls-cert-file"
	TLSKeyFileFlagName        = "tls-key-file"
	ConfigShortFlagName       = "c"
	ConfigLongFlagName        = "config"
	TrustedSubnetFlagName     = "t"
)

// Ошибки разора флагов.
var (
	ErrInvalidFlagValue                = errors.New("invalid flag value")
	ErrInvalidDefaultServerAddress     = errors.New("invalid default server address")
	ErrInvalidDefaultGRPCServerAddress = errors.New("invalid default GRPC server address")
)

// flagsConfig - настройки получаемые из флагов.
type flagsConfig struct {
	ServerAddress     *string
	GRPCServerAddress *string
	BaseURL           *string
	EnableLogs        *bool
	LogLevel          *string
	LogFile           *string
	FileStoragePath   *string
	DatabaseDSN       *string
	AuditFile         *string
	AuditURL          *string
	EnableHTTPS       *bool
	TLSCertFile       *string
	TLSKeyFile        *string
	ConfigFile        *string
	TrustedSubnet     *string
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
		if flagsConf.GRPCServerAddress != nil {
			config.GRPCServerAddress = *flagsConf.GRPCServerAddress
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
		if flagsConf.TLSCertFile != nil {
			config.TLSCertFile = *flagsConf.TLSCertFile
		}
		if flagsConf.TLSKeyFile != nil {
			config.TLSKeyFile = *flagsConf.TLSKeyFile
		}
		if flagsConf.ConfigFile != nil {
			config.ConfigFile = flagsConf.ConfigFile
		}
		if flagsConf.TrustedSubnet != nil {
			config.TrustedSubnet = *flagsConf.TrustedSubnet
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

	grpcServerAddress := address{}
	err = grpcServerAddress.Set(DefaultGRPCServerAddress)
	if err != nil {
		return errors.Join(err, ErrInvalidDefaultGRPCServerAddress)
	}

	fs.Var(&serverAddress, ServerAddressFlagName, "address and port to run server")
	fs.Var(&grpcServerAddress, GRPCServerAddressFlagName, "address and port to run GRPC server")
	baseURL := fs.String(BaseURLFlagName, DefaultBaseURL, "address and port for output short url")
	enableLogs := fs.Bool(EnableLogsFlagName, DefaultEnableLogs, "enable output logs")
	logLevel := fs.String(LogLevelFlagName, DefaultLogLevel, "log level")
	logFile := fs.String(LogFileFlagName, "", "log file")
	fileStoragePath := fs.String(FileStoragePathFlagName, "", "file storage path")
	databaseDSN := fs.String(DatabaseDSNFlagName, "", "database DSN")
	auditFile := fs.String(AuditFileFlagName, "", "audit file")
	auditURL := fs.String(AuditURLFlagName, "", "audit URL")
	enableHTTPS := fs.Bool(EnableHTTPSFlagName, DefaultEnableHTTPS, "enable HTTPS on web-server")
	tlsCertFile := fs.String(TLSCertFileFlagName, DefaultTLSCertFile, "TLS certificate file")
	tlsKeyFile := fs.String(TLSKeyFileFlagName, DefaultTLSKeyFile, "TLS certificate key file")
	var configFile string
	fs.StringVar(&configFile, ConfigShortFlagName, "", "config json file")
	fs.StringVar(&configFile, ConfigLongFlagName, "", "config json file")
	trustedSubnet := fs.String(TrustedSubnetFlagName, "", "CIDR of trusted subnet")

	err = fs.Parse(os.Args[1:])
	if err != nil {
		return fmt.Errorf("error parsing config flags: %w", err)
	}

	if config == nil {
		return nil
	}

	if justIfSet {
		isSetFlagMap := map[string]bool{
			ServerAddressFlagName:     false,
			GRPCServerAddressFlagName: false,
			BaseURLFlagName:           false,
			EnableLogsFlagName:        false,
			LogLevelFlagName:          false,
			LogFileFlagName:           false,
			FileStoragePathFlagName:   false,
			DatabaseDSNFlagName:       false,
			AuditFileFlagName:         false,
			AuditURLFlagName:          false,
			EnableHTTPSFlagName:       false,
			TLSCertFileFlagName:       false,
			TLSKeyFileFlagName:        false,
			ConfigShortFlagName:       false,
			ConfigLongFlagName:        false,
			TrustedSubnetFlagName:     false,
		}

		fs.Visit(func(f *flag.Flag) {
			isSetFlagMap[f.Name] = true
		})

		if isSet, ok := isSetFlagMap[ServerAddressFlagName]; ok && isSet {
			addr := serverAddress.String()
			config.ServerAddress = &addr
		}
		if isSet, ok := isSetFlagMap[GRPCServerAddressFlagName]; ok && isSet {
			addr := grpcServerAddress.String()
			config.GRPCServerAddress = &addr
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
		if isSet, ok := isSetFlagMap[TLSCertFileFlagName]; ok && isSet {
			if tlsCertFile != nil && *tlsCertFile != "" {
				config.TLSCertFile = tlsCertFile
			}
		}
		if isSet, ok := isSetFlagMap[TLSKeyFileFlagName]; ok && isSet {
			if tlsKeyFile != nil && *tlsKeyFile != "" {
				config.TLSKeyFile = tlsKeyFile
			}
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
		if isSet, ok := isSetFlagMap[TrustedSubnetFlagName]; ok && isSet {
			if trustedSubnet != nil && *trustedSubnet != "" {
				config.TrustedSubnet = trustedSubnet
			}
		}
	} else {
		addr := serverAddress.String()
		config.ServerAddress = &addr
		grpcAddr := grpcServerAddress.String()
		config.GRPCServerAddress = &grpcAddr
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
		if tlsCertFile != nil && *tlsCertFile != "" {
			config.TLSCertFile = tlsCertFile
		}
		if tlsKeyFile != nil && *tlsKeyFile != "" {
			config.TLSKeyFile = tlsKeyFile
		}
		if configFile != "" {
			config.ConfigFile = &configFile
		}
		if trustedSubnet != nil && *trustedSubnet != "" {
			config.TrustedSubnet = trustedSubnet
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

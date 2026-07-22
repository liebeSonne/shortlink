package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseFlags(t *testing.T) {
	appLog1 := "app.log"
	fileStoragePath1 := "./file/path"
	databaseDSN1 := "host=localhost user=username password=123 dbname=db sslmode=disable"
	auditFile1 := "audit1.log"
	auditURL1 := "https://audit.url1.com"
	configFile1 := "config1.json"
	configFile2 := "config2.json"
	tlsCertFile1 := "./some/path/to/cert.pem"
	tlsKeyFile1 := "./some/path/to/key.pem"
	trustedSubnet1 := "192.168.1.0/24"

	var defaultCfg = Config{
		BaseURL:           DefaultBaseURL,
		ServerAddress:     DefaultServerAddress,
		GRPCServerAddress: DefaultGRPCServerAddress,
		EnableLogs:        DefaultEnableLogs,
		LogLevel:          DefaultLogLevel,
		EnableHTTPS:       DefaultEnableHTTPS,
		TLSKeyFile:        DefaultTLSKeyFile,
		TLSCertFile:       DefaultTLSCertFile,
	}

	type want struct {
		cfg Config
		err error
	}
	testCases := []struct {
		name string
		args []string
		want want
	}{
		{"default args", []string{}, want{defaultCfg, nil}},
		{"set -a flag", []string{"-a", "10.0.0.1:8000"}, want{makeModConfig(defaultCfg, func(c *Config) { c.ServerAddress = "10.0.0.1:8000" }), nil}},
		{"set --a flag", []string{"--a", "10.0.0.1:8000"}, want{makeModConfig(defaultCfg, func(c *Config) { c.ServerAddress = "10.0.0.1:8000" }), nil}},
		{"set --a flag with empty address", []string{"--a", ":8000"}, want{makeModConfig(defaultCfg, func(c *Config) { c.ServerAddress = ":8000" }), nil}},
		{"set -g flag", []string{"-g", "10.0.0.1:3333"}, want{makeModConfig(defaultCfg, func(c *Config) { c.GRPCServerAddress = "10.0.0.1:3333" }), nil}},
		{"set --g flag", []string{"--g", "10.0.0.1:3333"}, want{makeModConfig(defaultCfg, func(c *Config) { c.GRPCServerAddress = "10.0.0.1:3333" }), nil}},
		{"set --g flag with empty address", []string{"--g", ":3333"}, want{makeModConfig(defaultCfg, func(c *Config) { c.GRPCServerAddress = ":3333" }), nil}},
		{"set -b flag", []string{"-b", "10.0.0.1:8000"}, want{makeModConfig(defaultCfg, func(c *Config) { c.BaseURL = "10.0.0.1:8000" }), nil}},
		{"set --b flag", []string{"--b", "10.0.0.1:8000"}, want{makeModConfig(defaultCfg, func(c *Config) { c.BaseURL = "10.0.0.1:8000" }), nil}},
		{"set -b flag with schema", []string{"-b", "http://10.0.0.1:8000"}, want{makeModConfig(defaultCfg, func(c *Config) { c.BaseURL = "http://10.0.0.1:8000" }), nil}},
		{"set -b flag with string", []string{"-b", "some-string"}, want{makeModConfig(defaultCfg, func(c *Config) { c.BaseURL = "some-string" }), nil}},
		{"set -a and -b flag", []string{"-a", "10.0.0.2:8081", "-b", "http://127.0.0.2:8082"}, want{makeModConfig(defaultCfg, func(c *Config) { c.ServerAddress = "10.0.0.2:8081"; c.BaseURL = "http://127.0.0.2:8082" }), nil}},
		{"set -b and -a flag", []string{"-b", "http://127.0.0.2:8082", "-a", "10.0.0.2:8081"}, want{makeModConfig(defaultCfg, func(c *Config) { c.ServerAddress = "10.0.0.2:8081"; c.BaseURL = "http://127.0.0.2:8082" }), nil}},
		{"set -a flag with invalid value", []string{"-a", "invalid value"}, want{err: ErrInvalidFlagValue}},
		{"set -a flag with invalid format", []string{"-a", "10.0.0.1:8080:abc"}, want{err: ErrInvalidFlagValue}},
		{"set -a flag with empty port", []string{"-a", "10.0.0.1:"}, want{err: ErrInvalidFlagValue}},
		{"set -a flag with invalid port", []string{"-a", "10.0.0.1:abc"}, want{err: ErrInvalidFlagValue}},
		{"set -g flag with invalid value", []string{"-g", "invalid value"}, want{err: ErrInvalidFlagValue}},
		{"set -g flag with invalid format", []string{"-g", "10.0.0.1:3333:abc"}, want{err: ErrInvalidFlagValue}},
		{"set -g flag with empty port", []string{"-g", "10.0.0.1:"}, want{err: ErrInvalidFlagValue}},
		{"set -g flag with invalid port", []string{"-g", "10.0.0.1:abc"}, want{err: ErrInvalidFlagValue}},
		{"set -l flag", []string{"-l=true"}, want{makeModConfig(defaultCfg, func(c *Config) { c.EnableLogs = true }), nil}},
		{"set -ll flag", []string{"-ll", "error"}, want{makeModConfig(defaultCfg, func(c *Config) { c.LogLevel = "error" }), nil}},
		{"set --ll flag", []string{"--ll", "error"}, want{makeModConfig(defaultCfg, func(c *Config) { c.LogLevel = "error" }), nil}},
		{"set --ll flag empty", []string{"--ll", ""}, want{makeModConfig(defaultCfg, func(c *Config) { c.LogLevel = "" }), nil}},
		{"set --ll flag custom value", []string{"--ll", "custom value"}, want{makeModConfig(defaultCfg, func(c *Config) { c.LogLevel = "custom value" }), nil}},
		{"set -lf flag", []string{"-lf", appLog1}, want{makeModConfig(defaultCfg, func(c *Config) { c.LogFile = &appLog1 }), nil}},
		{"set -lf flag empty", []string{"-lf", ""}, want{makeModConfig(defaultCfg, func(c *Config) { c.LogLevel = DefaultLogLevel }), nil}},
		{"set --lf flag", []string{"--lf", appLog1}, want{makeModConfig(defaultCfg, func(c *Config) { c.LogFile = &appLog1 }), nil}},
		{"set -f flag", []string{"-f", fileStoragePath1}, want{makeModConfig(defaultCfg, func(c *Config) { c.FileStoragePath = &fileStoragePath1 }), nil}},
		{"set -f flag empty", []string{"-f", ""}, want{defaultCfg, nil}},
		{"set --f flag", []string{"--f", fileStoragePath1}, want{makeModConfig(defaultCfg, func(c *Config) { c.FileStoragePath = &fileStoragePath1 }), nil}},
		{"set -d flag", []string{"-d", databaseDSN1}, want{makeModConfig(defaultCfg, func(c *Config) { c.DatabaseDSN = &databaseDSN1 }), nil}},
		{"set -d flag empty", []string{"-d", ""}, want{makeModConfig(defaultCfg, func(c *Config) { c.DatabaseDSN = nil }), nil}},
		{"set --d flag", []string{"--d", databaseDSN1}, want{makeModConfig(defaultCfg, func(c *Config) { c.DatabaseDSN = &databaseDSN1 }), nil}},
		{"set -audit-file flag", []string{"-audit-file", auditFile1}, want{makeModConfig(defaultCfg, func(c *Config) { c.AuditFile = &auditFile1 }), nil}},
		{"set --audit-file flag", []string{"--audit-file", auditFile1}, want{makeModConfig(defaultCfg, func(c *Config) { c.AuditFile = &auditFile1 }), nil}},
		{"set --audit-file flag empty", []string{"--audit-file", ""}, want{makeModConfig(defaultCfg, func(c *Config) { c.AuditFile = nil }), nil}},
		{"set -audit-url flag", []string{"-audit-url", auditURL1}, want{makeModConfig(defaultCfg, func(c *Config) { c.AuditURL = &auditURL1 }), nil}},
		{"set --audit-url flag", []string{"-audit-url", auditURL1}, want{makeModConfig(defaultCfg, func(c *Config) { c.AuditURL = &auditURL1 }), nil}},
		{"set --audit-url flag empty", []string{"-audit-url", ""}, want{makeModConfig(defaultCfg, func(c *Config) { c.AuditURL = nil }), nil}},
		{"set -s flag=true", []string{"-s=true"}, want{makeModConfig(defaultCfg, func(c *Config) { c.EnableHTTPS = true }), nil}},
		{"set --s flag=true", []string{"--s=true"}, want{makeModConfig(defaultCfg, func(c *Config) { c.EnableHTTPS = true }), nil}},
		{"set -s flag=false", []string{"-s=false"}, want{makeModConfig(defaultCfg, func(c *Config) { c.EnableHTTPS = false }), nil}},
		{"set --s flag=false", []string{"--s=false"}, want{makeModConfig(defaultCfg, func(c *Config) { c.EnableHTTPS = false }), nil}},
		{"set -tls-cert-file", []string{"-tls-cert-file", tlsCertFile1}, want{makeModConfig(defaultCfg, func(c *Config) { c.TLSCertFile = tlsCertFile1 }), nil}},
		{"set -tls-key-file", []string{"-tls-key-file", tlsKeyFile1}, want{makeModConfig(defaultCfg, func(c *Config) { c.TLSKeyFile = tlsKeyFile1 }), nil}},
		{"set -c flag empty", []string{"-c", ""}, want{defaultCfg, nil}},
		{"set --c flag empty", []string{"--c", ""}, want{defaultCfg, nil}},
		{"set -config flag empty", []string{"-config", ""}, want{defaultCfg, nil}},
		{"set --config flag empty", []string{"--config", ""}, want{defaultCfg, nil}},
		{"set -c flag", []string{"-c", configFile1}, want{makeModConfig(defaultCfg, func(c *Config) { c.ConfigFile = &configFile1 }), nil}},
		{"set --c flag", []string{"--c", configFile1}, want{makeModConfig(defaultCfg, func(c *Config) { c.ConfigFile = &configFile1 }), nil}},
		{"set -config flag", []string{"-config", configFile1}, want{makeModConfig(defaultCfg, func(c *Config) { c.ConfigFile = &configFile1 }), nil}},
		{"set --config flag", []string{"--config", configFile1}, want{makeModConfig(defaultCfg, func(c *Config) { c.ConfigFile = &configFile1 }), nil}},
		{"set -c flag and -config flag", []string{"-c", configFile1, "-config", configFile2}, want{makeModConfig(defaultCfg, func(c *Config) { c.ConfigFile = &configFile2 }), nil}},
		{"set -config flag and -c flag", []string{"-config", configFile1, "-c", configFile2}, want{makeModConfig(defaultCfg, func(c *Config) { c.ConfigFile = &configFile2 }), nil}},
		{"set -c flag and empty -config flag", []string{"-c", configFile1, "-config", ""}, want{defaultCfg, nil}},
		{"set -config flag and empty -c flag", []string{"-config", configFile1, "-c", ""}, want{defaultCfg, nil}},
		{"set empty -c flag and -config flag", []string{"-c", "", "-config", configFile2}, want{makeModConfig(defaultCfg, func(c *Config) { c.ConfigFile = &configFile2 }), nil}},
		{"set empty -config flag and -c flag", []string{"-config", "", "-c", configFile2}, want{makeModConfig(defaultCfg, func(c *Config) { c.ConfigFile = &configFile2 }), nil}},
		{"set -t flag", []string{"-t", trustedSubnet1}, want{makeModConfig(defaultCfg, func(c *Config) { c.TrustedSubnet = trustedSubnet1 }), nil}},
		{"set --t flag", []string{"-t", trustedSubnet1}, want{makeModConfig(defaultCfg, func(c *Config) { c.TrustedSubnet = trustedSubnet1 }), nil}},
		{"set --t flag empty", []string{"-t", ""}, want{makeModConfig(defaultCfg, func(c *Config) { c.TrustedSubnet = "" }), nil}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			appID := "test"
			oldArgs := os.Args
			args := []string{appID}
			args = append(args, tc.args...)
			os.Args = args
			defer func() {
				os.Args = oldArgs
			}()

			conf := Config{}
			err := parseFlags(appID, &conf)

			if tc.want.err != nil {
				require.Error(t, err)
				assert.ErrorContains(t, err, tc.want.err.Error())
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.want.cfg, conf)
		})
	}
}

func TestParseFlagsConfig(t *testing.T) {
	serverAddress1 := "10.10.10.10:1111"
	grpcServerAddress1 := "10.10.10.10:3333"
	baseURL1 := "http://127.0.0.1:2222"
	enableLogsTrue := true
	enableHTTPSTrue := true
	logLevel1 := LogLevelError
	appLog1 := "app.log"
	fileStoragePath1 := "./file/path"
	databaseDSN1 := "host=localhost user=username password=123 dbname=db sslmode=disable"
	configFile1 := "config1.json"
	tlsCertFile1 := "./some/path/to/cert.pem"
	tlsKeyFile1 := "./some/path/to/key.pem"
	trustedSubnet1 := "192.168.1.0/24"

	defaultServerAddress := DefaultServerAddress
	defaultGRPCServerAddress := DefaultGRPCServerAddress
	defaultBaseURL := DefaultBaseURL
	defaultEnableLogs := DefaultEnableLogs
	defaultLogLevel := DefaultLogLevel
	defaultTLSCertFile1 := DefaultTLSCertFile
	defaultTLSKeyFile1 := DefaultTLSKeyFile

	defaultFlagConfig := flagsConfig{
		ServerAddress:     &defaultServerAddress,
		GRPCServerAddress: &defaultGRPCServerAddress,
		BaseURL:           &defaultBaseURL,
		EnableLogs:        &defaultEnableLogs,
		LogLevel:          &defaultLogLevel,
		EnableHTTPS:       &defaultEnableLogs,
		TLSCertFile:       &defaultTLSCertFile1,
		TLSKeyFile:        &defaultTLSKeyFile1,
	}

	type when struct {
		args []string
	}
	type on struct {
		justIfSet bool
	}
	type want struct {
		cfg flagsConfig
		err error
	}
	testCases := []struct {
		name string
		when when
		on   on
		want want
	}{
		// and just if set
		{
			"empty args and just if set",
			when{[]string{}},
			on{true},
			want{flagsConfig{}, nil},
		},
		{
			"set -a flag and just if set",
			when{[]string{"-a", serverAddress1}},
			on{true},
			want{flagsConfig{ServerAddress: &serverAddress1}, nil},
		},
		{
			"set -g flag and just if set",
			when{[]string{"-g", grpcServerAddress1}},
			on{true},
			want{flagsConfig{GRPCServerAddress: &grpcServerAddress1}, nil},
		},
		{
			"set -b flag and just if set",
			when{[]string{"-b", baseURL1}},
			on{true},
			want{flagsConfig{BaseURL: &baseURL1}, nil},
		},
		{
			"set -l flag and just if set",
			when{[]string{"-l=true"}},
			on{true},
			want{flagsConfig{EnableLogs: &enableLogsTrue}, nil},
		},
		{
			"set -l without value flag and just if set",
			when{[]string{"-l"}},
			on{true},
			want{flagsConfig{EnableLogs: &enableLogsTrue}, nil},
		},
		{
			"set -ll flag and just if set",
			when{[]string{"-ll", logLevel1}},
			on{true},
			want{flagsConfig{LogLevel: &logLevel1}, nil},
		},
		{
			"set -lf flag and just if set",
			when{[]string{"-lf", appLog1}},
			on{true},
			want{flagsConfig{LogFile: &appLog1}, nil},
		},
		{
			"set -lf flag empty and just if set",
			when{[]string{"-lf", ""}},
			on{true},
			want{flagsConfig{LogFile: nil}, nil},
		},
		{
			"set -f flag and just if set",
			when{[]string{"-f", fileStoragePath1}},
			on{true},
			want{flagsConfig{FileStoragePath: &fileStoragePath1}, nil},
		},
		{
			"set -f flag empty and just if set",
			when{[]string{"-f", ""}},
			on{true},
			want{flagsConfig{FileStoragePath: nil}, nil},
		},
		{
			"set -d flag and just if set",
			when{[]string{"-d", databaseDSN1}},
			on{true},
			want{flagsConfig{DatabaseDSN: &databaseDSN1}, nil},
		},
		{
			"set -d flag empty and just if set",
			when{[]string{"-d", ""}},
			on{true},
			want{flagsConfig{DatabaseDSN: nil}, nil},
		},
		{
			"set -a flag with invalid value and just if set",
			when{[]string{"-a", "invalid value"}},
			on{true},
			want{err: ErrInvalidFlagValue},
		},
		{
			"set -a -b -l -ll flags and just if set",
			when{[]string{"-a", serverAddress1, "-b", baseURL1, "-l=true", "-ll", logLevel1}},
			on{true},
			want{flagsConfig{ServerAddress: &serverAddress1, BaseURL: &baseURL1, EnableLogs: &enableLogsTrue, LogLevel: &logLevel1}, nil},
		},
		{
			"set -a flag with invalid format and just if set",
			when{[]string{"-a", "10.0.0.1:8080:abc"}},
			on{true},
			want{err: ErrInvalidFlagValue},
		},
		{
			"set -a flag with empty port and just if set",
			when{[]string{"-a", "10.0.0.1:"}},
			on{true},
			want{err: ErrInvalidFlagValue},
		},
		{
			"set -a flag with invalid port and just if set",
			when{[]string{"-a", "10.0.0.1:abc"}},
			on{true},
			want{err: ErrInvalidFlagValue},
		},
		{
			"set -g flag with invalid value and just if set",
			when{[]string{"-g", "invalid value"}},
			on{true},
			want{err: ErrInvalidFlagValue},
		},
		{
			"set -g flag with invalid format and just if set",
			when{[]string{"-g", "10.0.0.1:3333:abc"}},
			on{true},
			want{err: ErrInvalidFlagValue},
		},
		{
			"set -g flag with empty port and just if set",
			when{[]string{"-g", "10.0.0.1:"}},
			on{true},
			want{err: ErrInvalidFlagValue},
		},
		{
			"set -g flag with invalid port and just if set",
			when{[]string{"-g", "10.0.0.1:abc"}},
			on{true},
			want{err: ErrInvalidFlagValue},
		},
		{
			"set -s flag and just if set",
			when{[]string{"-s=true"}},
			on{true},
			want{flagsConfig{EnableHTTPS: &enableHTTPSTrue}, nil},
		},
		{
			"set -s without value flag and just if set",
			when{[]string{"-s"}},
			on{true},
			want{flagsConfig{EnableHTTPS: &enableHTTPSTrue}, nil},
		},
		{
			"set -tls-cert-file flag and just if set",
			when{[]string{"-tls-cert-file", tlsCertFile1}},
			on{true},
			want{flagsConfig{TLSCertFile: &tlsCertFile1}, nil},
		},
		{
			"set -tls-key-file flag and just if set",
			when{[]string{"-tls-key-file", tlsKeyFile1}},
			on{true},
			want{flagsConfig{TLSKeyFile: &tlsKeyFile1}, nil},
		},
		{
			"set -c flag and just if set",
			when{[]string{"-c", configFile1}},
			on{true},
			want{flagsConfig{ConfigFile: &configFile1}, nil},
		},
		{
			"set -config flag and just if set",
			when{[]string{"-config", configFile1}},
			on{true},
			want{flagsConfig{ConfigFile: &configFile1}, nil},
		},
		{
			"set -t without value flag and just if set",
			when{[]string{"-t", trustedSubnet1}},
			on{true},
			want{flagsConfig{TrustedSubnet: &trustedSubnet1}, nil},
		},
		// and not just if set
		{
			"empty args and not just if set",
			when{[]string{}},
			on{false},
			want{defaultFlagConfig, nil},
		},
		{
			"set -a flag and not just if set",
			when{[]string{"-a", serverAddress1}},
			on{false},
			want{makeModFlagConfig(defaultFlagConfig, func(c *flagsConfig) { c.ServerAddress = &serverAddress1 }), nil},
		},
		{
			"set -g flag and not just if set",
			when{[]string{"-g", grpcServerAddress1}},
			on{false},
			want{makeModFlagConfig(defaultFlagConfig, func(c *flagsConfig) { c.GRPCServerAddress = &grpcServerAddress1 }), nil},
		},
		{
			"set -b flag and not just if set",
			when{[]string{"-b", baseURL1}},
			on{false},
			want{makeModFlagConfig(defaultFlagConfig, func(c *flagsConfig) { c.BaseURL = &baseURL1 }), nil},
		},
		{
			"set -l flag and not just if set",
			when{[]string{"-l=true"}},
			on{false},
			want{makeModFlagConfig(defaultFlagConfig, func(c *flagsConfig) { c.EnableLogs = &enableLogsTrue }), nil},
		},
		{
			"set -lf flag and not just if set",
			when{[]string{"-lf", appLog1}},
			on{false},
			want{makeModFlagConfig(defaultFlagConfig, func(c *flagsConfig) { c.LogFile = &appLog1 }), nil},
		},
		{
			"set -lf flag empty and not just if set",
			when{[]string{"-lf", ""}},
			on{false},
			want{makeModFlagConfig(defaultFlagConfig, func(c *flagsConfig) { c.LogFile = nil }), nil},
		},
		{
			"set -f flag and not just if set",
			when{[]string{"-f", fileStoragePath1}},
			on{false},
			want{makeModFlagConfig(defaultFlagConfig, func(c *flagsConfig) { c.FileStoragePath = &fileStoragePath1 }), nil},
		},
		{
			"set -f flag empty and not just if set",
			when{[]string{"-f", ""}},
			on{false},
			want{makeModFlagConfig(defaultFlagConfig, func(c *flagsConfig) { c.FileStoragePath = nil }), nil},
		},
		{
			"set -d flag and not just if set",
			when{[]string{"-d", databaseDSN1}},
			on{false},
			want{makeModFlagConfig(defaultFlagConfig, func(c *flagsConfig) { c.DatabaseDSN = &databaseDSN1 }), nil},
		},
		{
			"set -d flag empty and not just if set",
			when{[]string{"-d", ""}},
			on{false},
			want{makeModFlagConfig(defaultFlagConfig, func(c *flagsConfig) { c.DatabaseDSN = nil }), nil},
		},
		{
			"set -a -b -l flags and not just if set",
			when{[]string{"-a", serverAddress1, "-b", baseURL1, "-l=true"}},
			on{false},
			want{makeModFlagConfig(defaultFlagConfig, func(c *flagsConfig) {
				c.ServerAddress = &serverAddress1
				c.BaseURL = &baseURL1
				c.EnableLogs = &enableLogsTrue
			}), nil},
		},
		{
			"set -s flag and not just if set",
			when{[]string{"-s=true"}},
			on{false},
			want{makeModFlagConfig(defaultFlagConfig, func(c *flagsConfig) { c.EnableHTTPS = &enableHTTPSTrue }), nil},
		},
		{
			"set -tls-cert-file flag and not just if set",
			when{[]string{"-tls-cert-file", tlsCertFile1}},
			on{false},
			want{makeModFlagConfig(defaultFlagConfig, func(c *flagsConfig) { c.TLSCertFile = &tlsCertFile1 }), nil},
		},
		{
			"set -tls-key-file flag and not just if set",
			when{[]string{"-tls-key-file", tlsKeyFile1}},
			on{false},
			want{makeModFlagConfig(defaultFlagConfig, func(c *flagsConfig) { c.TLSKeyFile = &tlsKeyFile1 }), nil},
		},
		{
			"set -a flag with invalid value and not just if set",
			when{[]string{"-a", "invalid value"}},
			on{false},
			want{err: ErrInvalidFlagValue},
		},
		{
			"set -a flag with invalid format and not just if set",
			when{[]string{"-a", "10.0.0.1:8080:abc"}},
			on{false},
			want{err: ErrInvalidFlagValue},
		},
		{
			"set -a flag with empty port and not just if set",
			when{[]string{"-a", "10.0.0.1:"}},
			on{false},
			want{err: ErrInvalidFlagValue},
		},
		{
			"set -a flag with invalid port and not just if set",
			when{[]string{"-a", "10.0.0.1:abc"}},
			on{false},
			want{err: ErrInvalidFlagValue},
		},
		{
			"set -g flag with invalid value and not just if set",
			when{[]string{"-g", "invalid value"}},
			on{false},
			want{err: ErrInvalidFlagValue},
		},
		{
			"set -g flag with invalid format and not just if set",
			when{[]string{"-g", "10.0.0.1:8080:abc"}},
			on{false},
			want{err: ErrInvalidFlagValue},
		},
		{
			"set -g flag with empty port and not just if set",
			when{[]string{"-g", "10.0.0.1:"}},
			on{false},
			want{err: ErrInvalidFlagValue},
		},
		{
			"set -g flag with invalid port and not just if set",
			when{[]string{"-g", "10.0.0.1:abc"}},
			on{false},
			want{err: ErrInvalidFlagValue},
		},
		{
			"set -c flag and not just if set",
			when{[]string{"-c", configFile1}},
			on{false},
			want{makeModFlagConfig(defaultFlagConfig, func(c *flagsConfig) { c.ConfigFile = &configFile1 }), nil},
		},
		{
			"set -c flag empty and not just if set",
			when{[]string{"-c", ""}},
			on{false},
			want{makeModFlagConfig(defaultFlagConfig, func(c *flagsConfig) { c.ConfigFile = nil }), nil},
		},
		{
			"set -config flag and not just if set",
			when{[]string{"-config", configFile1}},
			on{false},
			want{makeModFlagConfig(defaultFlagConfig, func(c *flagsConfig) { c.ConfigFile = &configFile1 }), nil},
		},
		{
			"set -config flag empty and not just if set",
			when{[]string{"-config", ""}},
			on{false},
			want{makeModFlagConfig(defaultFlagConfig, func(c *flagsConfig) { c.ConfigFile = nil }), nil},
		},
		{
			"set -t flag and not just if set",
			when{[]string{"-t", trustedSubnet1}},
			on{false},
			want{makeModFlagConfig(defaultFlagConfig, func(c *flagsConfig) { c.TrustedSubnet = &trustedSubnet1 }), nil},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			appID := "test"
			oldArgs := os.Args
			args := []string{appID}
			args = append(args, tc.when.args...)
			os.Args = args
			defer func() {
				os.Args = oldArgs
			}()

			cfg := flagsConfig{}
			err := parseFlagsConfig(appID, &cfg, tc.on.justIfSet)

			if tc.want.err != nil {
				require.Error(t, err)
				assert.ErrorContains(t, err, tc.want.err.Error())
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.want.cfg, cfg)
		})
	}
}

func TestAddress_Set(t *testing.T) {
	type on struct {
		value string
	}
	type want struct {
		addr address
		err  error
	}
	testCases := []struct {
		name string
		on   on
		want want
	}{
		{"empty value", on{""}, want{err: ErrInvalidFlagValue}},
		{"invalid value", on{"invalid value"}, want{err: ErrInvalidFlagValue}},
		{"invalid format", on{"127.0.0.1:8888:abc"}, want{err: ErrInvalidFlagValue}},
		{"empty host and empty port", on{":"}, want{err: ErrInvalidFlagValue}},
		{"host and empty port", on{"127.0.0.1:"}, want{err: ErrInvalidFlagValue}},
		{"empty host and port", on{":8888"}, want{address{"", 8888}, nil}},
		{"host and port", on{"127.0.0.1:8888"}, want{address{"127.0.0.1", 8888}, nil}},
		{"string host and port", on{"string:8888"}, want{address{"string", 8888}, nil}},
		{"invalid port", on{"127.0.0.1:abc"}, want{err: ErrInvalidFlagValue}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			a := address{}

			err := a.Set(tc.on.value)

			if tc.want.err != nil {
				require.Error(t, err)
				assert.ErrorContains(t, err, tc.want.err.Error())
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.want.addr, a)
		})
	}
}

func TestAddress_String(t *testing.T) {
	testCases := []struct {
		name string
		on   address
		want string
	}{
		{"empty", address{}, ":0"},
		{"empty host and port", address{"", 8888}, ":8888"},
		{"host and empty port", address{Host: ""}, ":0"},
		{"host and port", address{"127.0.0.1", 8888}, "127.0.0.1:8888"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.on.String()
			require.Equal(t, tc.want, result)
		})
	}
}

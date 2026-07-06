package config

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestParseEnv(t *testing.T) {
	appLog := "app.log"
	fileStoragePath := "file/path"
	databaseDSN := "host=localhost user=user password=password dbname=db sslmode=disabled"
	authCookieTokenKey1 := "token-key-1"
	authTokenExpiresStr1 := "10m"
	authTokenExpiresDuration1 := time.Minute * 10
	authSecretKey1 := "secret123"
	auditFile1 := "audit1.log"
	auditURL1 := "https://audit.url1.com"
	configFile1 := "config1.json"
	tlsCertFile1 := "./some/path/to/cert.pem"
	tlsKeyFile1 := "./some/path/to/key.pem"

	type on struct {
		prefix string
	}
	type when struct {
		envs map[string]string
	}
	type want struct {
		conf Config
		err  error
	}
	testCases := []struct {
		name string
		on   on
		when when
		want want
	}{
		{
			"default",
			on{""},
			when{map[string]string{}},
			want{defaultConfig, nil},
		},
		{
			"server address",
			on{""},
			when{map[string]string{
				getEnvNameWithPrefix("", ServerAddressEnvName): "127.0.0.1:8888",
			}},
			want{MakeModConfig(defaultConfig, func(c *Config) {
				c.ServerAddress = "127.0.0.1:8888"
			}), nil},
		},
		{
			"base url",
			on{""},
			when{map[string]string{
				getEnvNameWithPrefix("", BaseURLEnvName): "http://127.0.0.1:8888",
			}},
			want{
				MakeModConfig(defaultConfig, func(c *Config) {
					c.BaseURL = "http://127.0.0.1:8888"
				}), nil},
		},
		{
			"enable logs",
			on{""},
			when{map[string]string{
				getEnvNameWithPrefix("", EnableLogsEnvName): "true",
			}},
			want{
				MakeModConfig(defaultConfig, func(c *Config) {
					c.EnableLogs = true
				}), nil},
		},
		{
			"not enable logs",
			on{""},
			when{map[string]string{
				getEnvNameWithPrefix("", EnableLogsEnvName): "false",
			}},
			want{
				MakeModConfig(defaultConfig, func(c *Config) {
					c.EnableLogs = false
				}), nil},
		},
		{
			"log level",
			on{""},
			when{map[string]string{
				getEnvNameWithPrefix("", LogLevelEnvName): LogLevelError,
			}},
			want{
				MakeModConfig(defaultConfig, func(c *Config) {
					c.LogLevel = LogLevelError
				}), nil},
		},
		{
			"log file",
			on{""},
			when{map[string]string{
				getEnvNameWithPrefix("", LogFileEnvName): appLog,
			}},
			want{
				MakeModConfig(defaultConfig, func(c *Config) {
					c.LogFile = &appLog
				}), nil},
		},
		{
			"file storage path",
			on{""},
			when{map[string]string{
				getEnvNameWithPrefix("", FileStoragePathEnvName): fileStoragePath,
			}},
			want{
				MakeModConfig(defaultConfig, func(c *Config) {
					c.FileStoragePath = &fileStoragePath
				}), nil},
		},
		{
			"database dsn",
			on{""},
			when{map[string]string{
				getEnvNameWithPrefix("", DatabaseDSNEnvName): databaseDSN,
			}},
			want{
				MakeModConfig(defaultConfig, func(c *Config) {
					c.DatabaseDSN = &databaseDSN
				}), nil},
		},
		{
			"audit file",
			on{""},
			when{map[string]string{
				getEnvNameWithPrefix("", AuditFileEnvName): auditFile1,
			}},
			want{
				MakeModConfig(defaultConfig, func(c *Config) {
					c.AuditFile = &auditFile1
				}), nil},
		},
		{
			"audit url",
			on{""},
			when{map[string]string{
				getEnvNameWithPrefix("", AuditURLEnvName): auditURL1,
			}},
			want{
				MakeModConfig(defaultConfig, func(c *Config) {
					c.AuditURL = &auditURL1
				}), nil},
		},
		{
			"enable HTTPS",
			on{""},
			when{map[string]string{
				getEnvNameWithPrefix("", EnableHTTPSEnvName): "true",
			}},
			want{
				MakeModConfig(defaultConfig, func(c *Config) {
					c.EnableHTTPS = true
				}), nil},
		},
		{
			"not enable HTTPS",
			on{""},
			when{map[string]string{
				getEnvNameWithPrefix("", EnableHTTPSEnvName): "false",
			}},
			want{
				MakeModConfig(defaultConfig, func(c *Config) {
					c.EnableHTTPS = false
				}), nil},
		},
		{
			"TLS cert file",
			on{""},
			when{map[string]string{
				getEnvNameWithPrefix("", TLSCertFileEnvName): tlsCertFile1,
			}},
			want{
				MakeModConfig(defaultConfig, func(c *Config) {
					c.TLSCertFile = tlsCertFile1
				}), nil},
		},
		{
			"TLS key file",
			on{""},
			when{map[string]string{
				getEnvNameWithPrefix("", TLSKeyFileEnvName): tlsKeyFile1,
			}},
			want{
				MakeModConfig(defaultConfig, func(c *Config) {
					c.TLSKeyFile = tlsKeyFile1
				}), nil},
		},
		{
			"config file",
			on{""},
			when{map[string]string{
				getEnvNameWithPrefix("", ConfigFileEnvName): configFile1,
			}},
			want{
				MakeModConfig(defaultConfig, func(c *Config) {
					c.ConfigFile = &configFile1
				}), nil},
		},
		{
			"all env",
			on{""},
			when{map[string]string{
				getEnvNameWithPrefix("", ServerAddressEnvName):      "127.0.0.1:8888",
				getEnvNameWithPrefix("", BaseURLEnvName):            "http://127.0.0.2:8000",
				getEnvNameWithPrefix("", EnableLogsEnvName):         "true",
				getEnvNameWithPrefix("", LogLevelEnvName):           LogLevelError,
				getEnvNameWithPrefix("", LogFileEnvName):            appLog,
				getEnvNameWithPrefix("", FileStoragePathEnvName):    fileStoragePath,
				getEnvNameWithPrefix("", DatabaseDSNEnvName):        databaseDSN,
				getEnvNameWithPrefix("", AuthCookieTokenKeyEnvName): authCookieTokenKey1,
				getEnvNameWithPrefix("", AuthTokenExpiresEnvName):   authTokenExpiresStr1,
				getEnvNameWithPrefix("", AuthSecretKeyEnvName):      authSecretKey1,
				getEnvNameWithPrefix("", AuditFileEnvName):          auditFile1,
				getEnvNameWithPrefix("", AuditURLEnvName):           auditURL1,
				getEnvNameWithPrefix("", EnableHTTPSEnvName):        "true",
				getEnvNameWithPrefix("", TLSCertFileEnvName):        tlsCertFile1,
				getEnvNameWithPrefix("", TLSKeyFileEnvName):         tlsKeyFile1,
				getEnvNameWithPrefix("", ConfigFileEnvName):         configFile1,
			}},
			want{Config{
				ServerAddress:      "127.0.0.1:8888",
				BaseURL:            "http://127.0.0.2:8000",
				EnableLogs:         true,
				LogLevel:           LogLevelError,
				LogFile:            &appLog,
				FileStoragePath:    &fileStoragePath,
				DatabaseDSN:        &databaseDSN,
				AuthCookieTokenKey: authCookieTokenKey1,
				AuthTokenExpires:   authTokenExpiresDuration1,
				AuthSecretKey:      authSecretKey1,
				AuditFile:          &auditFile1,
				AuditURL:           &auditURL1,
				EnableHTTPS:        true,
				TLSCertFile:        tlsCertFile1,
				TLSKeyFile:         tlsKeyFile1,
				ConfigFile:         &configFile1,
			}, nil},
		},
		{
			"with prefix",
			on{"app_id"},
			when{map[string]string{
				getEnvNameWithPrefix("APP_ID", ServerAddressEnvName):      "127.0.0.1:8888",
				getEnvNameWithPrefix("APP_ID", BaseURLEnvName):            "http://127.0.0.2:8000",
				getEnvNameWithPrefix("APP_ID", EnableLogsEnvName):         "true",
				getEnvNameWithPrefix("APP_ID", LogLevelEnvName):           LogLevelError,
				getEnvNameWithPrefix("APP_ID", LogFileEnvName):            appLog,
				getEnvNameWithPrefix("APP_ID", FileStoragePathEnvName):    fileStoragePath,
				getEnvNameWithPrefix("APP_ID", DatabaseDSNEnvName):        databaseDSN,
				getEnvNameWithPrefix("APP_ID", AuthCookieTokenKeyEnvName): authCookieTokenKey1,
				getEnvNameWithPrefix("APP_ID", AuthTokenExpiresEnvName):   authTokenExpiresStr1,
				getEnvNameWithPrefix("APP_ID", AuthSecretKeyEnvName):      authSecretKey1,
				getEnvNameWithPrefix("APP_ID", AuditFileEnvName):          auditFile1,
				getEnvNameWithPrefix("APP_ID", AuditURLEnvName):           auditURL1,
				getEnvNameWithPrefix("APP_ID", EnableHTTPSEnvName):        "true",
				getEnvNameWithPrefix("APP_ID", TLSCertFileEnvName):        tlsCertFile1,
				getEnvNameWithPrefix("APP_ID", TLSKeyFileEnvName):         tlsKeyFile1,
				getEnvNameWithPrefix("APP_ID", ConfigFileEnvName):         configFile1,
			}},
			want{Config{
				ServerAddress:      "127.0.0.1:8888",
				BaseURL:            "http://127.0.0.2:8000",
				EnableLogs:         true,
				LogLevel:           LogLevelError,
				LogFile:            &appLog,
				FileStoragePath:    &fileStoragePath,
				DatabaseDSN:        &databaseDSN,
				AuthCookieTokenKey: authCookieTokenKey1,
				AuthTokenExpires:   authTokenExpiresDuration1,
				AuthSecretKey:      authSecretKey1,
				AuditFile:          &auditFile1,
				AuditURL:           &auditURL1,
				EnableHTTPS:        true,
				TLSCertFile:        tlsCertFile1,
				TLSKeyFile:         tlsKeyFile1,
				ConfigFile:         &configFile1,
			}, nil},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			oldEnv := os.Environ()
			os.Clearenv()
			for k, v := range tc.when.envs {
				t.Setenv(k, v)
			}
			t.Cleanup(func() {
				os.Clearenv()
				for _, pair := range oldEnv {
					kv := strings.SplitN(pair, "=", 2)
					_ = os.Setenv(kv[0], kv[1])
				}
			})

			conf := Config{}
			envPrefix := getEnvNameWithPrefix(tc.on.prefix, "")
			err := ParseEnv(envPrefix, &conf)

			if tc.want.err != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tc.want.err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.want.conf, conf)
		})
	}
}

func TestGetEnvNameWithPrefix(t *testing.T) {
	type on struct {
		prefix  string
		envName string
	}
	testCases := []struct {
		name string
		on   on
		want string
	}{
		{"prefix", on{"prefix", "env1_name"}, "PREFIX_ENV1_NAME"},
		{"empty prefix", on{"", "env1_name"}, "ENV1_NAME"},
		{"empty env name", on{"prefix", ""}, "PREFIX_"},
		{"empty prefix and empty env name", on{"", ""}, ""},
		{"number in prefix", on{"1", "env1_name"}, "1_ENV1_NAME"},
		{"number env name", on{"prefix", "1"}, "PREFIX_1"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := getEnvNameWithPrefix(tc.on.prefix, tc.on.envName)
			require.Equal(t, tc.want, result)
		})
	}
}

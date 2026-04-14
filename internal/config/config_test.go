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
			want{Config{
				ServerAddress:      DefaultServerAddress,
				BaseURL:            DefaultBaseURL,
				EnableLogs:         DefaultEnableLogs,
				LogLevel:           DefaultLogLevel,
				AuthCookieTokenKey: DefaultAuthCookieTokenKey,
				AuthTokenExpires:   DefaultAuthTokenExpires,
				AuthSecretKey:      DefaultAuthSecretKey,
			}, nil},
		},
		{
			"server address",
			on{""},
			when{map[string]string{
				getEnvNameWithPrefix("", ServerAddressEnvName): "127.0.0.1:8888",
			}},
			want{Config{
				ServerAddress:      "127.0.0.1:8888",
				BaseURL:            DefaultBaseURL,
				EnableLogs:         DefaultEnableLogs,
				LogLevel:           DefaultLogLevel,
				AuthCookieTokenKey: DefaultAuthCookieTokenKey,
				AuthTokenExpires:   DefaultAuthTokenExpires,
				AuthSecretKey:      DefaultAuthSecretKey,
			}, nil},
		},
		{
			"base url",
			on{""},
			when{map[string]string{
				getEnvNameWithPrefix("", BaseURLEnvName): "http://127.0.0.1:8888",
			}},
			want{Config{
				ServerAddress:      DefaultServerAddress,
				BaseURL:            "http://127.0.0.1:8888",
				EnableLogs:         DefaultEnableLogs,
				LogLevel:           DefaultLogLevel,
				AuthCookieTokenKey: DefaultAuthCookieTokenKey,
				AuthTokenExpires:   DefaultAuthTokenExpires,
				AuthSecretKey:      DefaultAuthSecretKey,
			}, nil},
		},
		{
			"enable logs",
			on{""},
			when{map[string]string{
				getEnvNameWithPrefix("", EnableLogsEnvName): "true",
			}},
			want{Config{
				ServerAddress:      DefaultServerAddress,
				BaseURL:            DefaultBaseURL,
				EnableLogs:         true,
				LogLevel:           DefaultLogLevel,
				AuthCookieTokenKey: DefaultAuthCookieTokenKey,
				AuthTokenExpires:   DefaultAuthTokenExpires,
				AuthSecretKey:      DefaultAuthSecretKey,
			}, nil},
		},
		{
			"not enable logs",
			on{""},
			when{map[string]string{
				getEnvNameWithPrefix("", EnableLogsEnvName): "false",
			}},
			want{Config{
				ServerAddress:      DefaultServerAddress,
				BaseURL:            DefaultBaseURL,
				EnableLogs:         false,
				LogLevel:           DefaultLogLevel,
				AuthCookieTokenKey: DefaultAuthCookieTokenKey,
				AuthTokenExpires:   DefaultAuthTokenExpires,
				AuthSecretKey:      DefaultAuthSecretKey,
			}, nil},
		},
		{
			"log level",
			on{""},
			when{map[string]string{
				getEnvNameWithPrefix("", LogLevelEnvName): LogLevelError,
			}},
			want{Config{
				ServerAddress:      DefaultServerAddress,
				BaseURL:            DefaultBaseURL,
				EnableLogs:         DefaultEnableLogs,
				LogLevel:           LogLevelError,
				AuthCookieTokenKey: DefaultAuthCookieTokenKey,
				AuthTokenExpires:   DefaultAuthTokenExpires,
				AuthSecretKey:      DefaultAuthSecretKey,
			}, nil},
		},
		{
			"log file",
			on{""},
			when{map[string]string{
				getEnvNameWithPrefix("", LogFileEnvName): appLog,
			}},
			want{Config{
				ServerAddress:      DefaultServerAddress,
				BaseURL:            DefaultBaseURL,
				EnableLogs:         DefaultEnableLogs,
				LogLevel:           DefaultLogLevel,
				LogFile:            &appLog,
				AuthCookieTokenKey: DefaultAuthCookieTokenKey,
				AuthTokenExpires:   DefaultAuthTokenExpires,
				AuthSecretKey:      DefaultAuthSecretKey,
			}, nil},
		},
		{
			"file storage path",
			on{""},
			when{map[string]string{
				getEnvNameWithPrefix("", FileStoragePathEnvName): fileStoragePath,
			}},
			want{Config{
				ServerAddress:      DefaultServerAddress,
				BaseURL:            DefaultBaseURL,
				EnableLogs:         DefaultEnableLogs,
				LogLevel:           DefaultLogLevel,
				FileStoragePath:    &fileStoragePath,
				AuthCookieTokenKey: DefaultAuthCookieTokenKey,
				AuthTokenExpires:   DefaultAuthTokenExpires,
				AuthSecretKey:      DefaultAuthSecretKey,
			}, nil},
		},
		{
			"database dsn",
			on{""},
			when{map[string]string{
				getEnvNameWithPrefix("", DatabaseDSNEnvName): databaseDSN,
			}},
			want{Config{
				ServerAddress:      DefaultServerAddress,
				BaseURL:            DefaultBaseURL,
				EnableLogs:         DefaultEnableLogs,
				LogLevel:           DefaultLogLevel,
				DatabaseDSN:        &databaseDSN,
				AuthCookieTokenKey: DefaultAuthCookieTokenKey,
				AuthTokenExpires:   DefaultAuthTokenExpires,
				AuthSecretKey:      DefaultAuthSecretKey,
			}, nil},
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

package config

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseEnv(t *testing.T) {
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
				ServerAddress: DefaultServerAddress,
				BaseURL:       DefaultBaseURL,
				EnableLogs:    DefaultEnableLogs,
				LogLevel:      DefaultLogLevel,
			}, nil},
		},
		{
			"server address",
			on{""},
			when{map[string]string{
				getEnvNameWithPrefix("", ServerAddressEnvName): "127.0.0.1:8888",
			}},
			want{Config{
				ServerAddress: "127.0.0.1:8888",
				BaseURL:       DefaultBaseURL,
				EnableLogs:    DefaultEnableLogs,
				LogLevel:      DefaultLogLevel,
			}, nil},
		},
		{
			"base url",
			on{""},
			when{map[string]string{
				getEnvNameWithPrefix("", BaseURLEnvName): "http://127.0.0.1:8888",
			}},
			want{Config{
				ServerAddress: DefaultServerAddress,
				BaseURL:       "http://127.0.0.1:8888",
				EnableLogs:    DefaultEnableLogs,
				LogLevel:      DefaultLogLevel,
			}, nil},
		},
		{
			"enable logs",
			on{""},
			when{map[string]string{
				getEnvNameWithPrefix("", EnableLogsEnvName): "true",
			}},
			want{Config{
				ServerAddress: DefaultServerAddress,
				BaseURL:       DefaultBaseURL,
				EnableLogs:    true,
				LogLevel:      DefaultLogLevel,
			}, nil},
		},
		{
			"not enable logs",
			on{""},
			when{map[string]string{
				getEnvNameWithPrefix("", EnableLogsEnvName): "false",
			}},
			want{Config{
				ServerAddress: DefaultServerAddress,
				BaseURL:       DefaultBaseURL,
				EnableLogs:    false,
				LogLevel:      DefaultLogLevel,
			}, nil},
		},
		{
			"log level",
			on{""},
			when{map[string]string{
				getEnvNameWithPrefix("", LogLevelEnvName): LogLevelError,
			}},
			want{Config{
				ServerAddress: DefaultServerAddress,
				BaseURL:       DefaultBaseURL,
				EnableLogs:    DefaultEnableLogs,
				LogLevel:      LogLevelError,
			}, nil},
		},
		{
			"all env",
			on{""},
			when{map[string]string{
				getEnvNameWithPrefix("", ServerAddressEnvName): "127.0.0.1:8888",
				getEnvNameWithPrefix("", BaseURLEnvName):       "http://127.0.0.2:8000",
				getEnvNameWithPrefix("", EnableLogsEnvName):    "true",
				getEnvNameWithPrefix("", LogLevelEnvName):      LogLevelError,
			}},
			want{Config{
				ServerAddress: "127.0.0.1:8888",
				BaseURL:       "http://127.0.0.2:8000",
				EnableLogs:    true,
				LogLevel:      LogLevelError,
			}, nil},
		},
		{
			"with prefix",
			on{"app_id"},
			when{map[string]string{
				getEnvNameWithPrefix("APP_ID", ServerAddressEnvName): "127.0.0.1:8888",
				getEnvNameWithPrefix("APP_ID", BaseURLEnvName):       "http://127.0.0.2:8000",
				getEnvNameWithPrefix("APP_ID", EnableLogsEnvName):    "true",
				getEnvNameWithPrefix("APP_ID", LogLevelEnvName):      LogLevelError,
			}},
			want{Config{
				ServerAddress: "127.0.0.1:8888",
				BaseURL:       "http://127.0.0.2:8000",
				EnableLogs:    true,
				LogLevel:      LogLevelError,
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

package config

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	type when struct {
		appID     string
		envPrefix string
		args      []string
		envs      map[string]string
	}
	type want struct {
		conf Config
		err  error
	}
	testCases := []struct {
		name string
		when when
		want want
	}{
		{
			"default",
			when{},
			want{
				Config{
					ServerAddress: DefaultServerAddress,
					BaseURL:       DefaultBaseURL,
					EnableLogs:    DefaultEnableLogs,
					LogLevel:      DefaultLogLevel,
				},
				nil,
			},
		},
		{
			"from env with prefix",
			when{
				"app_id", "prefix",
				[]string{},
				map[string]string{
					getEnvNameWithPrefix("prefix", ServerAddressEnvName): "127.0.0.1:8787",
					getEnvNameWithPrefix("prefix", BaseURLEnvName):       "http://127.0.0.2:8888",
					getEnvNameWithPrefix("prefix", EnableLogsEnvName):    "true",
					getEnvNameWithPrefix("prefix", LogLevelEnvName):      "error",
				},
			},
			want{
				Config{
					ServerAddress: "127.0.0.1:8787",
					BaseURL:       "http://127.0.0.2:8888",
					EnableLogs:    true,
					LogLevel:      "error",
				},
				nil,
			},
		},
		{
			"from env without prefix",
			when{
				"app_id", "",
				[]string{},
				map[string]string{
					getEnvNameWithPrefix("", ServerAddressEnvName): "127.0.0.1:8787",
					getEnvNameWithPrefix("", BaseURLEnvName):       "http://127.0.0.2:8888",
					getEnvNameWithPrefix("", EnableLogsEnvName):    "true",
					getEnvNameWithPrefix("", LogLevelEnvName):      "error",
				},
			},
			want{
				Config{
					ServerAddress: "127.0.0.1:8787",
					BaseURL:       "http://127.0.0.2:8888",
					EnableLogs:    true,
					LogLevel:      "error",
				},
				nil,
			},
		},
		{
			"from flags",
			when{
				"", "",
				[]string{"-a", "127.0.0.1:8787", "-b", "http://127.0.0.2:8888", "-l=true", "-ll", "error"},
				map[string]string{},
			},
			want{
				Config{
					ServerAddress: "127.0.0.1:8787",
					BaseURL:       "http://127.0.0.2:8888",
					EnableLogs:    true,
					LogLevel:      "error",
				},
				nil,
			},
		},
		{
			"server address from env and base url from flags",
			when{
				"", "",
				[]string{"-a", "127.0.0.1:8787", "-b", "http://127.0.0.2:8888"},
				map[string]string{
					getEnvNameWithPrefix("", ServerAddressEnvName): "127.0.0.10:7777",
					getEnvNameWithPrefix("", EnableLogsEnvName):    "true",
					getEnvNameWithPrefix("", LogLevelEnvName):      "error",
				},
			},
			want{
				Config{
					ServerAddress: "127.0.0.10:7777",
					BaseURL:       "http://127.0.0.2:8888",
					EnableLogs:    true,
					LogLevel:      "error",
				},
				nil,
			},
		},
		{
			"server address from flags and base url from env",
			when{
				"", "",
				[]string{"-a", "127.0.0.1:8787", "-b", "http://127.0.0.2:8888"},
				map[string]string{
					getEnvNameWithPrefix("", BaseURLEnvName):    "http://127.0.0.2:8888",
					getEnvNameWithPrefix("", EnableLogsEnvName): "true",
					getEnvNameWithPrefix("", LogLevelEnvName):   "error",
				},
			},
			want{
				Config{
					ServerAddress: "127.0.0.1:8787",
					BaseURL:       "http://127.0.0.2:8888",
					EnableLogs:    true,
					LogLevel:      "error",
				},
				nil,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			oldArgs := os.Args
			oldEnv := os.Environ()

			args := []string{""}
			args = append(args, tc.when.args...)
			os.Args = args

			os.Clearenv()
			for k, v := range tc.when.envs {
				t.Setenv(k, v)
			}
			t.Cleanup(func() {
				os.Args = oldArgs
				os.Clearenv()
				for _, pair := range oldEnv {
					kv := strings.SplitN(pair, "=", 2)
					_ = os.Setenv(kv[0], kv[1])
				}
			})

			conf, err := LoadConfig(tc.when.appID, tc.when.envPrefix)

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

func TestMergeFlagsConfig(t *testing.T) {
	serverAddress1 := "10.10.10.10:1111"
	baseURL1 := "http://127.0.0.1:2222"
	enableLogsTrue := true
	logLevel1 := LogLevelError

	type on struct {
		fCfg     flagsConfig
		envNames []string
	}
	type want struct {
		cfg Config
	}
	testCases := []struct {
		name string
		on
		want want
	}{
		{
			"empty env names",
			on{flagsConfig{&serverAddress1, &baseURL1, &enableLogsTrue, &logLevel1}, []string{}},
			want{Config{}},
		},
		{
			"server address env name",
			on{flagsConfig{&serverAddress1, &baseURL1, &enableLogsTrue, &logLevel1}, []string{ServerAddressEnvName}},
			want{Config{ServerAddress: serverAddress1}},
		},
		{
			"base url env name",
			on{flagsConfig{&serverAddress1, &baseURL1, &enableLogsTrue, &logLevel1}, []string{BaseURLEnvName}},
			want{Config{BaseURL: baseURL1}},
		},
		{
			"enable logs env name",
			on{flagsConfig{&serverAddress1, &baseURL1, &enableLogsTrue, &logLevel1}, []string{EnableLogsEnvName}},
			want{Config{EnableLogs: enableLogsTrue}},
		},
		{
			"unknown env name",
			on{flagsConfig{&serverAddress1, &baseURL1, &enableLogsTrue, &logLevel1}, []string{"unknown"}},
			want{Config{}},
		},
		{
			"server address and base url and enable logs env names",
			on{flagsConfig{&serverAddress1, &baseURL1, &enableLogsTrue, &logLevel1}, []string{ServerAddressEnvName, BaseURLEnvName, EnableLogsEnvName}},
			want{Config{ServerAddress: serverAddress1, BaseURL: baseURL1, EnableLogs: enableLogsTrue}},
		},
		{
			"log level env name",
			on{flagsConfig{&serverAddress1, &baseURL1, &enableLogsTrue, &logLevel1}, []string{LogLevelEnvName}},
			want{Config{LogLevel: logLevel1}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := Config{}
			mergeFlagsConfig(tc.on.fCfg, &cfg, tc.on.envNames)

			require.Equal(t, tc.want.cfg, cfg)
		})
	}
}

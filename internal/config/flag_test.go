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

	type want struct {
		cfg Config
		err error
	}
	testCases := []struct {
		name string
		args []string
		want want
	}{
		{"default args", []string{}, want{Config{DefaultServerAddress, DefaultBaseURL, DefaultEnableLogs, DefaultLogLevel, nil, nil}, nil}},
		{"set -a flag", []string{"-a", "10.0.0.1:8000"}, want{Config{"10.0.0.1:8000", DefaultBaseURL, DefaultEnableLogs, DefaultLogLevel, nil, nil}, nil}},
		{"set --a flag", []string{"--a", "10.0.0.1:8000"}, want{Config{"10.0.0.1:8000", DefaultBaseURL, DefaultEnableLogs, DefaultLogLevel, nil, nil}, nil}},
		{"set -b flag", []string{"-b", "10.0.0.1:8000"}, want{Config{DefaultServerAddress, "10.0.0.1:8000", DefaultEnableLogs, DefaultLogLevel, nil, nil}, nil}},
		{"set --a flag with empty address", []string{"--a", ":8000"}, want{Config{":8000", DefaultBaseURL, DefaultEnableLogs, DefaultLogLevel, nil, nil}, nil}},
		{"set --b flag", []string{"--b", "10.0.0.1:8000"}, want{Config{DefaultServerAddress, "10.0.0.1:8000", DefaultEnableLogs, DefaultLogLevel, nil, nil}, nil}},
		{"set -b flag with schema", []string{"-b", "http://10.0.0.1:8000"}, want{Config{DefaultServerAddress, "http://10.0.0.1:8000", DefaultEnableLogs, DefaultLogLevel, nil, nil}, nil}},
		{"set -b flag with string", []string{"-b", "some-string"}, want{Config{DefaultServerAddress, "some-string", DefaultEnableLogs, DefaultLogLevel, nil, nil}, nil}},
		{"set -a and -b flag", []string{"-a", "10.0.0.2:8081", "-b", "http://127.0.0.2:8082"}, want{Config{"10.0.0.2:8081", "http://127.0.0.2:8082", DefaultEnableLogs, DefaultLogLevel, nil, nil}, nil}},
		{"set -b and -a flag", []string{"-b", "http://127.0.0.2:8082", "-a", "10.0.0.2:8081"}, want{Config{"10.0.0.2:8081", "http://127.0.0.2:8082", DefaultEnableLogs, DefaultLogLevel, nil, nil}, nil}},
		{"set -a flag with invalid value", []string{"-a", "invalid value"}, want{err: ErrInvalidFlagValue}},
		{"set -a flag with invalid format", []string{"-a", "10.0.0.1:8080:abc"}, want{err: ErrInvalidFlagValue}},
		{"set -a flag with empty port", []string{"-a", "10.0.0.1:"}, want{err: ErrInvalidFlagValue}},
		{"set -a flag with invalid port", []string{"-a", "10.0.0.1:abc"}, want{err: ErrInvalidFlagValue}},
		{"set -l flag", []string{"-l=true"}, want{Config{DefaultServerAddress, DefaultBaseURL, true, DefaultLogLevel, nil, nil}, nil}},
		{"set -ll flag", []string{"-ll", "error"}, want{Config{DefaultServerAddress, DefaultBaseURL, DefaultEnableLogs, "error", nil, nil}, nil}},
		{"set --ll flag", []string{"--ll", "error"}, want{Config{DefaultServerAddress, DefaultBaseURL, DefaultEnableLogs, "error", nil, nil}, nil}},
		{"set --ll flag empty", []string{"--ll", ""}, want{Config{DefaultServerAddress, DefaultBaseURL, DefaultEnableLogs, "", nil, nil}, nil}},
		{"set --ll flag custom value", []string{"--ll", "custom value"}, want{Config{DefaultServerAddress, DefaultBaseURL, DefaultEnableLogs, "custom value", nil, nil}, nil}},
		{"set -lf flag", []string{"-lf", appLog1}, want{Config{DefaultServerAddress, DefaultBaseURL, DefaultEnableLogs, DefaultLogLevel, &appLog1, nil}, nil}},
		{"set -lf flag empty", []string{"-lf", ""}, want{Config{DefaultServerAddress, DefaultBaseURL, DefaultEnableLogs, DefaultLogLevel, nil, nil}, nil}},
		{"set --lf flag", []string{"--lf", appLog1}, want{Config{DefaultServerAddress, DefaultBaseURL, DefaultEnableLogs, DefaultLogLevel, &appLog1, nil}, nil}},
		{"set -f flag", []string{"-f", fileStoragePath1}, want{Config{DefaultServerAddress, DefaultBaseURL, DefaultEnableLogs, DefaultLogLevel, nil, &fileStoragePath1}, nil}},
		{"set -f flag empty", []string{"-f", ""}, want{Config{DefaultServerAddress, DefaultBaseURL, DefaultEnableLogs, DefaultLogLevel, nil, nil}, nil}},
		{"set --f flag", []string{"--f", fileStoragePath1}, want{Config{DefaultServerAddress, DefaultBaseURL, DefaultEnableLogs, DefaultLogLevel, nil, &fileStoragePath1}, nil}},
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
	baseURL1 := "http://127.0.0.1:2222"
	enableLogsTrue := true
	logLevel1 := LogLevelError
	appLog1 := "app.log"
	fileStoragePath1 := "./file/path"

	defaultServerAddress := DefaultServerAddress
	defaultBaseURL := DefaultBaseURL
	defaultEnableLogs := DefaultEnableLogs
	defaultLogLevel := DefaultLogLevel

	defaultFlagConfig := flagsConfig{
		ServerAddress: &defaultServerAddress,
		BaseURL:       &defaultBaseURL,
		EnableLogs:    &defaultEnableLogs,
		LogLevel:      &defaultLogLevel,
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
			want{flagsConfig{ServerAddress: &serverAddress1, BaseURL: &defaultBaseURL, EnableLogs: &defaultEnableLogs, LogLevel: &defaultLogLevel}, nil},
		},
		{
			"set -b flag and not just if set",
			when{[]string{"-b", baseURL1}},
			on{false},
			want{flagsConfig{ServerAddress: &defaultServerAddress, BaseURL: &baseURL1, EnableLogs: &defaultEnableLogs, LogLevel: &defaultLogLevel}, nil},
		},
		{
			"set -l flag and not just if set",
			when{[]string{"-l=true"}},
			on{false},
			want{flagsConfig{ServerAddress: &defaultServerAddress, BaseURL: &defaultBaseURL, EnableLogs: &enableLogsTrue, LogLevel: &defaultLogLevel}, nil},
		},
		{
			"set -lf flag and not just if set",
			when{[]string{"-lf", appLog1}},
			on{false},
			want{flagsConfig{ServerAddress: &defaultServerAddress, BaseURL: &defaultBaseURL, EnableLogs: &defaultEnableLogs, LogLevel: &defaultLogLevel, LogFile: &appLog1}, nil},
		},
		{
			"set -lf flag empty and not just if set",
			when{[]string{"-lf", ""}},
			on{false},
			want{flagsConfig{ServerAddress: &defaultServerAddress, BaseURL: &defaultBaseURL, EnableLogs: &defaultEnableLogs, LogLevel: &defaultLogLevel, LogFile: nil}, nil},
		},
		{
			"set -f flag and not just if set",
			when{[]string{"-f", fileStoragePath1}},
			on{false},
			want{flagsConfig{ServerAddress: &defaultServerAddress, BaseURL: &defaultBaseURL, EnableLogs: &defaultEnableLogs, LogLevel: &defaultLogLevel, FileStoragePath: &fileStoragePath1}, nil},
		},
		{
			"set -f flag empty and not just if set",
			when{[]string{"-f", ""}},
			on{false},
			want{flagsConfig{ServerAddress: &defaultServerAddress, BaseURL: &defaultBaseURL, EnableLogs: &defaultEnableLogs, LogLevel: &defaultLogLevel, FileStoragePath: nil}, nil},
		},
		{
			"set -a -b -l flags and not just if set",
			when{[]string{"-a", serverAddress1, "-b", baseURL1, "-l=true"}},
			on{false},
			want{flagsConfig{ServerAddress: &serverAddress1, BaseURL: &baseURL1, EnableLogs: &enableLogsTrue, LogLevel: &defaultLogLevel}, nil},
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

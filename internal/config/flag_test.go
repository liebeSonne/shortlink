package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseFlags(t *testing.T) {
	defaultConf := Config{
		ServerAddress: "10.10.10.10:1111",
		BaseURL:       "127.2.2.2:2222",
	}

	type want struct {
		serverAddress string
		baseURL       string
		err           error
	}
	testCases := []struct {
		name string
		args []string
		want want
	}{
		{"default args", []string{}, want{DefaultServerAddress, DefaultBaseURL, nil}},
		{"set -a flag", []string{"-a", "10.0.0.1:8000"}, want{"10.0.0.1:8000", DefaultBaseURL, nil}},
		{"set --a flag", []string{"--a", "10.0.0.1:8000"}, want{"10.0.0.1:8000", DefaultBaseURL, nil}},
		{"set -b flag", []string{"-b", "10.0.0.1:8000"}, want{DefaultServerAddress, "10.0.0.1:8000", nil}},
		{"set --a flag with empty address", []string{"--a", ":8000"}, want{":8000", DefaultBaseURL, nil}},
		{"set --b flag", []string{"--b", "10.0.0.1:8000"}, want{DefaultServerAddress, "10.0.0.1:8000", nil}},
		{"set -b flag with schema", []string{"-b", "http://10.0.0.1:8000"}, want{DefaultServerAddress, "http://10.0.0.1:8000", nil}},
		{"set -b flag with string", []string{"-b", "some-string"}, want{DefaultServerAddress, "some-string", nil}},
		{"set -a and -b flag", []string{"-a", "10.0.0.2:8081", "-b", "http://127.0.0.2:8082"}, want{"10.0.0.2:8081", "http://127.0.0.2:8082", nil}},
		{"set -b and -a flag", []string{"-b", "http://127.0.0.2:8082", "-a", "10.0.0.2:8081"}, want{"10.0.0.2:8081", "http://127.0.0.2:8082", nil}},
		{"set -a flag with invalid value", []string{"-a", "invalid value"}, want{err: ErrInvalidFlagValue}},
		{"set -a flag with invalid format", []string{"-a", "10.0.0.1:8080:abc"}, want{err: ErrInvalidFlagValue}},
		{"set -a flag with empty port", []string{"-a", "10.0.0.1:"}, want{err: ErrInvalidFlagValue}},
		{"set -a flag with invalid port", []string{"-a", "10.0.0.1:abc"}, want{err: ErrInvalidFlagValue}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			conf := defaultConf

			appID := "test"
			oldArgs := os.Args
			args := []string{appID}
			args = append(args, tc.args...)
			os.Args = args
			defer func() {
				os.Args = oldArgs
			}()

			err := parseFlags(appID, &conf)

			if tc.want.err != nil {
				require.Error(t, err)
				assert.ErrorContains(t, err, tc.want.err.Error())
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.want.serverAddress, conf.ServerAddress)
			assert.Equal(t, tc.want.baseURL, conf.BaseURL)
		})
	}
}

func TestParseFlagsConfig(t *testing.T) {
	serverAddress1 := "10.10.10.10:1111"
	baseURL1 := "http://127.0.0.1:2222"
	enableLogsTrue := true
	defaultServerAddress := DefaultServerAddress
	defaultBaseURL := DefaultBaseURL
	defaultEnableLogs := DefaultEnableLogs

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
			when{[]string{"-l", "true"}},
			on{true},
			want{flagsConfig{EnableLogs: &enableLogsTrue}, nil},
		},
		{
			"set -a -b -l flags and just if set",
			when{[]string{"-a", serverAddress1, "-b", baseURL1, "-l", "true"}},
			on{true},
			want{flagsConfig{ServerAddress: &serverAddress1, BaseURL: &baseURL1, EnableLogs: &enableLogsTrue}, nil},
		},
		{
			"set -a flag with invalid value and just if set",
			when{[]string{"-a", "invalid value"}},
			on{true},
			want{err: ErrInvalidFlagValue},
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
			want{flagsConfig{ServerAddress: &defaultServerAddress, BaseURL: &defaultBaseURL, EnableLogs: &defaultEnableLogs}, nil},
		},
		{
			"set -a flag and not just if set",
			when{[]string{"-a", serverAddress1}},
			on{false},
			want{flagsConfig{ServerAddress: &serverAddress1, BaseURL: &defaultBaseURL, EnableLogs: &defaultEnableLogs}, nil},
		},
		{
			"set -b flag and not just if set",
			when{[]string{"-b", baseURL1}},
			on{false},
			want{flagsConfig{ServerAddress: &defaultServerAddress, BaseURL: &baseURL1, EnableLogs: &defaultEnableLogs}, nil},
		},
		{
			"set -l flag and not just if set",
			when{[]string{"-l", "true"}},
			on{false},
			want{flagsConfig{ServerAddress: &defaultServerAddress, BaseURL: &defaultBaseURL, EnableLogs: &enableLogsTrue}, nil},
		},
		{
			"set -a -b -l flags and not just if set",
			when{[]string{"-a", serverAddress1, "-b", baseURL1, "-l", "true"}},
			on{false},
			want{flagsConfig{ServerAddress: &serverAddress1, BaseURL: &baseURL1, EnableLogs: &enableLogsTrue}, nil},
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

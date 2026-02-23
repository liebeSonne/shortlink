package main

import (
	"os"
	"testing"

	"github.com/liebeSonne/shortlink/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseFlags(t *testing.T) {
	defaultConf := config.Config{
		ServerAddress: "10.10.10.10:1111",
		URLAddress:    "127.2.2.2:2222",
	}

	type want struct {
		serverAddress string
		urlAddress    string
		err           error
	}
	testCases := []struct {
		name string
		args []string
		want want
	}{
		{"default args", []string{}, want{defaultServerAddress, defaultURLAddress, nil}},
		{"set -a flag", []string{"-a", "10.0.0.1:8000"}, want{"10.0.0.1:8000", defaultURLAddress, nil}},
		{"set --a flag", []string{"--a", "10.0.0.1:8000"}, want{"10.0.0.1:8000", defaultURLAddress, nil}},
		{"set -b flag", []string{"-b", "10.0.0.1:8000"}, want{defaultServerAddress, "10.0.0.1:8000", nil}},
		{"set --a flag with empty address", []string{"--a", ":8000"}, want{":8000", defaultURLAddress, nil}},
		{"set --b flag", []string{"--b", "10.0.0.1:8000"}, want{defaultServerAddress, "10.0.0.1:8000", nil}},
		{"set -b flag with schema", []string{"-b", "http://10.0.0.1:8000"}, want{defaultServerAddress, "http://10.0.0.1:8000", nil}},
		{"set -b flag with string", []string{"-b", "some-string"}, want{defaultServerAddress, "some-string", nil}},
		{"set -a and b flag", []string{"-a", "10.0.0.2:8081", "-b", "http://127.0.0.2:8082"}, want{"10.0.0.2:8081", "http://127.0.0.2:8082", nil}},
		{"set -b and a flag", []string{"-b", "http://127.0.0.2:8082", "-a", "10.0.0.2:8081"}, want{"10.0.0.2:8081", "http://127.0.0.2:8082", nil}},
		{"set -a flag with invalid value", []string{"-a", "invalid value"}, want{defaultConf.ServerAddress, defaultConf.URLAddress, ErrInvalidFlagValue}},
		{"set -a flag with invalid format", []string{"-a", "10.0.0.1:8080:abc"}, want{defaultConf.ServerAddress, defaultConf.URLAddress, ErrInvalidFlagValue}},
		{"set -a flag with empty port", []string{"-a", "10.0.0.1:"}, want{defaultConf.ServerAddress, defaultConf.URLAddress, ErrInvalidFlagValue}},
		{"set -a flag with invalid port", []string{"-a", "10.0.0.1:abc"}, want{defaultConf.ServerAddress, defaultConf.URLAddress, ErrInvalidFlagValue}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			conf := defaultConf

			args := []string{""}
			args = append(args, tc.args...)
			os.Args = args

			err := parseFlags(&conf)

			if tc.want.err != nil {
				require.Error(t, err)
				assert.ErrorContains(t, err, tc.want.err.Error())
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tc.want.serverAddress, conf.ServerAddress)
			assert.Equal(t, tc.want.urlAddress, conf.URLAddress)
		})
	}
}

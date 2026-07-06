package config

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	appLog1 := "app.log"
	fileStoragePath1 := "./file/path"
	databaseDSN1 := "host=localhost user=username password=password dbname=db sslmode=disable"
	authCookieTokenKey1 := "token-key-1"
	authTokenExpiresStr1 := "10m"
	authTokenExpiresDuration1 := time.Minute * 10
	authSecretKey1 := "secret123"
	auditFile1 := "audit1.log"
	auditURL1 := "https://audit.url1.com"
	tlsCertFile1 := "./some/path/to/cert.pem"
	tlsKeyFile1 := "./some/path/to/key.pem"

	// make empty config file
	tempFile2, err := os.CreateTemp("", "config_*.json")
	require.NoError(t, err)
	defer os.Remove(tempFile2.Name())
	configFileData2 := `{}`
	err = os.WriteFile(tempFile2.Name(), []byte(configFileData2), 0644)
	require.NoError(t, err)
	configFile1 := tempFile2.Name()

	// make not empty config file
	tempFile, err := os.CreateTemp("", "config_*.json")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())
	configFileServerAddress1 := "1.1.1.1:1111"
	configFileBaseURL1 := "http://1.1.1.1:2222"
	configFileDatabaseDSN1 := "file database dns"
	configFileStoragePath1 := "./config/file/storage/path"
	configFileData1 := fmt.Sprintf(`{
		"server_address": "%s",
		"base_url": "%s",
		"file_storage_path": "%s",
		"database_dsn": "%s" 
	}`, configFileServerAddress1, configFileBaseURL1, configFileStoragePath1, configFileDatabaseDSN1)
	err = os.WriteFile(tempFile.Name(), []byte(configFileData1), 0644)
	require.NoError(t, err)
	configFile2 := tempFile.Name()

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
			want{defaultConfig, nil},
		},
		{
			"from env with prefix",
			when{
				"app_id", "prefix",
				[]string{},
				map[string]string{
					getEnvNameWithPrefix("prefix", ServerAddressEnvName):      "127.0.0.1:8787",
					getEnvNameWithPrefix("prefix", BaseURLEnvName):            "http://127.0.0.2:8888",
					getEnvNameWithPrefix("prefix", EnableLogsEnvName):         "true",
					getEnvNameWithPrefix("prefix", LogLevelEnvName):           "error",
					getEnvNameWithPrefix("prefix", LogFileEnvName):            appLog1,
					getEnvNameWithPrefix("prefix", FileStoragePathEnvName):    fileStoragePath1,
					getEnvNameWithPrefix("prefix", DatabaseDSNEnvName):        databaseDSN1,
					getEnvNameWithPrefix("prefix", AuthCookieTokenKeyEnvName): authCookieTokenKey1,
					getEnvNameWithPrefix("prefix", AuthTokenExpiresEnvName):   authTokenExpiresStr1,
					getEnvNameWithPrefix("prefix", AuthSecretKeyEnvName):      authSecretKey1,
					getEnvNameWithPrefix("prefix", AuditFileEnvName):          auditFile1,
					getEnvNameWithPrefix("prefix", AuditURLEnvName):           auditURL1,
					getEnvNameWithPrefix("prefix", EnableHTTPSEnvName):        "true",
					getEnvNameWithPrefix("prefix", TLSCertFileEnvName):        tlsCertFile1,
					getEnvNameWithPrefix("prefix", TLSKeyFileEnvName):         tlsKeyFile1,
					getEnvNameWithPrefix("prefix", ConfigFileEnvName):         configFile1,
				},
			},
			want{
				Config{
					ServerAddress:      "127.0.0.1:8787",
					BaseURL:            "http://127.0.0.2:8888",
					EnableLogs:         true,
					LogLevel:           "error",
					LogFile:            &appLog1,
					FileStoragePath:    &fileStoragePath1,
					DatabaseDSN:        &databaseDSN1,
					AuthCookieTokenKey: authCookieTokenKey1,
					AuthTokenExpires:   authTokenExpiresDuration1,
					AuthSecretKey:      authSecretKey1,
					AuditFile:          &auditFile1,
					AuditURL:           &auditURL1,
					EnableHTTPS:        true,
					TLSCertFile:        tlsCertFile1,
					TLSKeyFile:         tlsKeyFile1,
					ConfigFile:         &configFile1,
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
					getEnvNameWithPrefix("", ServerAddressEnvName):      "127.0.0.1:8787",
					getEnvNameWithPrefix("", BaseURLEnvName):            "http://127.0.0.2:8888",
					getEnvNameWithPrefix("", EnableLogsEnvName):         "true",
					getEnvNameWithPrefix("", LogLevelEnvName):           "error",
					getEnvNameWithPrefix("", LogFileEnvName):            appLog1,
					getEnvNameWithPrefix("", FileStoragePathEnvName):    fileStoragePath1,
					getEnvNameWithPrefix("", DatabaseDSNEnvName):        databaseDSN1,
					getEnvNameWithPrefix("", AuthCookieTokenKeyEnvName): authCookieTokenKey1,
					getEnvNameWithPrefix("", AuthTokenExpiresEnvName):   authTokenExpiresStr1,
					getEnvNameWithPrefix("", AuthSecretKeyEnvName):      authSecretKey1,
					getEnvNameWithPrefix("", AuditFileEnvName):          auditFile1,
					getEnvNameWithPrefix("", AuditURLEnvName):           auditURL1,
					getEnvNameWithPrefix("", EnableHTTPSEnvName):        "true",
					getEnvNameWithPrefix("", TLSCertFileEnvName):        tlsCertFile1,
					getEnvNameWithPrefix("", TLSKeyFileEnvName):         tlsKeyFile1,
					getEnvNameWithPrefix("", ConfigFileEnvName):         configFile1,
				},
			},
			want{
				Config{
					ServerAddress:      "127.0.0.1:8787",
					BaseURL:            "http://127.0.0.2:8888",
					EnableLogs:         true,
					LogLevel:           "error",
					LogFile:            &appLog1,
					FileStoragePath:    &fileStoragePath1,
					DatabaseDSN:        &databaseDSN1,
					AuthCookieTokenKey: authCookieTokenKey1,
					AuthTokenExpires:   authTokenExpiresDuration1,
					AuthSecretKey:      authSecretKey1,
					AuditFile:          &auditFile1,
					AuditURL:           &auditURL1,
					EnableHTTPS:        true,
					TLSCertFile:        tlsCertFile1,
					TLSKeyFile:         tlsKeyFile1,
					ConfigFile:         &configFile1,
				},
				nil,
			},
		},
		{
			"from flags",
			when{
				"", "",
				[]string{
					"-a", "127.0.0.1:8787",
					"-b", "http://127.0.0.2:8888",
					"-l=true",
					"-ll", "error",
					"-lf", appLog1,
					"-f", fileStoragePath1,
					"-d", databaseDSN1,
					"--audit-file", auditFile1,
					"--audit-url", auditURL1,
					"-s=true",
					"-tls-cert-file", tlsCertFile1,
					"-tls-key-file", tlsKeyFile1,
					"-c", configFile1,
				},
				map[string]string{},
			},
			want{
				Config{
					ServerAddress:      "127.0.0.1:8787",
					BaseURL:            "http://127.0.0.2:8888",
					EnableLogs:         true,
					LogLevel:           "error",
					LogFile:            &appLog1,
					FileStoragePath:    &fileStoragePath1,
					DatabaseDSN:        &databaseDSN1,
					AuthCookieTokenKey: DefaultAuthCookieTokenKey,
					AuthTokenExpires:   DefaultAuthTokenExpires,
					AuthSecretKey:      DefaultAuthSecretKey,
					AuditFile:          &auditFile1,
					AuditURL:           &auditURL1,
					EnableHTTPS:        true,
					TLSCertFile:        tlsCertFile1,
					TLSKeyFile:         tlsKeyFile1,
					ConfigFile:         &configFile1,
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
					getEnvNameWithPrefix("", ServerAddressEnvName):      "127.0.0.10:7777",
					getEnvNameWithPrefix("", EnableLogsEnvName):         "true",
					getEnvNameWithPrefix("", LogLevelEnvName):           "error",
					getEnvNameWithPrefix("", LogFileEnvName):            appLog1,
					getEnvNameWithPrefix("", FileStoragePathEnvName):    fileStoragePath1,
					getEnvNameWithPrefix("", DatabaseDSNEnvName):        databaseDSN1,
					getEnvNameWithPrefix("", AuthCookieTokenKeyEnvName): authCookieTokenKey1,
					getEnvNameWithPrefix("", AuthTokenExpiresEnvName):   authTokenExpiresStr1,
					getEnvNameWithPrefix("", AuthSecretKeyEnvName):      authSecretKey1,
					getEnvNameWithPrefix("", EnableHTTPSEnvName):        "true",
					getEnvNameWithPrefix("", TLSCertFileEnvName):        tlsCertFile1,
					getEnvNameWithPrefix("", TLSKeyFileEnvName):         tlsKeyFile1,
					getEnvNameWithPrefix("", ConfigFileEnvName):         configFile1,
				},
			},
			want{
				Config{
					ServerAddress:      "127.0.0.10:7777",
					BaseURL:            "http://127.0.0.2:8888",
					EnableLogs:         true,
					LogLevel:           "error",
					LogFile:            &appLog1,
					FileStoragePath:    &fileStoragePath1,
					DatabaseDSN:        &databaseDSN1,
					AuthCookieTokenKey: authCookieTokenKey1,
					AuthTokenExpires:   authTokenExpiresDuration1,
					AuthSecretKey:      authSecretKey1,
					EnableHTTPS:        true,
					TLSCertFile:        tlsCertFile1,
					TLSKeyFile:         tlsKeyFile1,
					ConfigFile:         &configFile1,
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
					getEnvNameWithPrefix("", BaseURLEnvName):            "http://127.0.0.2:8888",
					getEnvNameWithPrefix("", EnableLogsEnvName):         "true",
					getEnvNameWithPrefix("", LogLevelEnvName):           "error",
					getEnvNameWithPrefix("", LogFileEnvName):            appLog1,
					getEnvNameWithPrefix("", FileStoragePathEnvName):    fileStoragePath1,
					getEnvNameWithPrefix("", DatabaseDSNEnvName):        databaseDSN1,
					getEnvNameWithPrefix("", AuthCookieTokenKeyEnvName): authCookieTokenKey1,
					getEnvNameWithPrefix("", AuthTokenExpiresEnvName):   authTokenExpiresStr1,
					getEnvNameWithPrefix("", AuthSecretKeyEnvName):      authSecretKey1,
					getEnvNameWithPrefix("", EnableHTTPSEnvName):        "true",
					getEnvNameWithPrefix("", TLSCertFileEnvName):        tlsCertFile1,
					getEnvNameWithPrefix("", TLSKeyFileEnvName):         tlsKeyFile1,
					getEnvNameWithPrefix("", ConfigFileEnvName):         configFile1,
				},
			},
			want{
				Config{
					ServerAddress:      "127.0.0.1:8787",
					BaseURL:            "http://127.0.0.2:8888",
					EnableLogs:         true,
					LogLevel:           "error",
					LogFile:            &appLog1,
					FileStoragePath:    &fileStoragePath1,
					DatabaseDSN:        &databaseDSN1,
					AuthCookieTokenKey: authCookieTokenKey1,
					AuthTokenExpires:   authTokenExpiresDuration1,
					AuthSecretKey:      authSecretKey1,
					EnableHTTPS:        true,
					TLSCertFile:        tlsCertFile1,
					TLSKeyFile:         tlsKeyFile1,
					ConfigFile:         &configFile1,
				},
				nil,
			},
		},
		{
			"from config file in env",
			when{
				"", "",
				[]string{},
				map[string]string{
					getEnvNameWithPrefix("", ConfigFileEnvName): configFile2,
				},
			},
			want{
				MakeModConfig(defaultConfig, func(c *Config) {
					c.ServerAddress = configFileServerAddress1
					c.BaseURL = configFileBaseURL1
					c.FileStoragePath = &configFileStoragePath1
					c.DatabaseDSN = &configFileDatabaseDSN1
					c.ConfigFile = &configFile2
				}),
				nil,
			},
		},
		{
			"from config file in flag",
			when{
				"", "",
				[]string{"-c", configFile2},
				map[string]string{},
			},
			want{
				MakeModConfig(defaultConfig, func(c *Config) {
					c.ServerAddress = configFileServerAddress1
					c.BaseURL = configFileBaseURL1
					c.FileStoragePath = &configFileStoragePath1
					c.DatabaseDSN = &configFileDatabaseDSN1
					c.ConfigFile = &configFile2
				}),
				nil,
			},
		},
		{
			"from config file with priority from flag and env",
			when{
				"", "",
				[]string{"-a", "127.0.0.1:8787", "-b", "http://127.0.0.2:8888"},
				map[string]string{
					getEnvNameWithPrefix("", BaseURLEnvName):    "http://127.0.0.2:8888",
					getEnvNameWithPrefix("", ConfigFileEnvName): configFile2,
				},
			},
			want{
				MakeModConfig(defaultConfig, func(c *Config) {
					c.ServerAddress = "127.0.0.1:8787"
					c.BaseURL = "http://127.0.0.2:8888"
					c.FileStoragePath = &configFileStoragePath1
					c.DatabaseDSN = &configFileDatabaseDSN1
					c.ConfigFile = &configFile2
				}),
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
	enableHTTPSTrue := true
	logLevel1 := LogLevelError
	logFile1 := "app.log"
	fileStoragePath1 := "./file/path"
	databaseDSN1 := "host=localhost user=username password=password dbname=db sslmode=disable"
	auditFile1 := "audit1.log"
	auditURL1 := "https://audit.url1.com"
	configFile1 := "config1.json"
	tlsCertFile1 := "./some/path/to/cert.pem"
	tlsKeyFile1 := "./some/path/to/key.pem"

	flagConfig1 := flagsConfig{
		&serverAddress1,
		&baseURL1,
		&enableLogsTrue,
		&logLevel1,
		&logFile1,
		&fileStoragePath1,
		&databaseDSN1,
		&auditFile1,
		&auditURL1,
		&enableHTTPSTrue,
		&tlsCertFile1,
		&tlsKeyFile1,
		&configFile1,
	}

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
			on{flagConfig1, []string{}},
			want{Config{}},
		},
		{
			"server address env name",
			on{flagConfig1, []string{ServerAddressEnvName}},
			want{Config{ServerAddress: serverAddress1}},
		},
		{
			"base url env name",
			on{flagConfig1, []string{BaseURLEnvName}},
			want{Config{BaseURL: baseURL1}},
		},
		{
			"enable logs env name",
			on{flagConfig1, []string{EnableLogsEnvName}},
			want{Config{EnableLogs: enableLogsTrue}},
		},
		{
			"unknown env name",
			on{flagConfig1, []string{"unknown"}},
			want{Config{}},
		},
		{
			"server address and base url and enable logs env names",
			on{flagConfig1, []string{ServerAddressEnvName, BaseURLEnvName, EnableLogsEnvName}},
			want{Config{ServerAddress: serverAddress1, BaseURL: baseURL1, EnableLogs: enableLogsTrue}},
		},
		{
			"log level env name",
			on{flagConfig1, []string{LogLevelEnvName}},
			want{Config{LogLevel: logLevel1}},
		},
		{
			"log file env name",
			on{flagConfig1, []string{LogFileEnvName}},
			want{Config{LogFile: &logFile1}},
		},
		{
			"file path storage env name",
			on{flagConfig1, []string{FileStoragePathEnvName}},
			want{Config{FileStoragePath: &fileStoragePath1}},
		},
		{
			"database DSN env name",
			on{flagConfig1, []string{DatabaseDSNEnvName}},
			want{Config{DatabaseDSN: &databaseDSN1}},
		},
		{
			"audit file env name",
			on{flagConfig1, []string{AuditFileEnvName}},
			want{Config{AuditFile: &auditFile1}},
		},
		{
			"audit url env name",
			on{flagConfig1, []string{AuditURLEnvName}},
			want{Config{AuditURL: &auditURL1}},
		},
		{
			"enable HTTPS env name",
			on{flagConfig1, []string{EnableHTTPSEnvName}},
			want{Config{EnableHTTPS: enableHTTPSTrue}},
		},
		{
			"TLS cert file env name",
			on{flagConfig1, []string{TLSCertFileEnvName}},
			want{Config{TLSCertFile: tlsCertFile1}},
		},
		{
			"TLS key file env name",
			on{flagConfig1, []string{TLSKeyFileEnvName}},
			want{Config{TLSKeyFile: tlsKeyFile1}},
		},
		{
			"config file env name",
			on{flagConfig1, []string{ConfigFileEnvName}},
			want{Config{ConfigFile: &configFile1}},
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

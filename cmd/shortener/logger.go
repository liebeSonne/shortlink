package main

import (
	"fmt"
	"io"
	"os"

	"github.com/liebeSonne/shortlink/internal/config"
	internalio "github.com/liebeSonne/shortlink/internal/io"
	applogger "github.com/liebeSonne/shortlink/internal/logger"
)

var configToLoggerLogLevelMap = map[string]applogger.LogLevel{
	config.LogLevelDebug: applogger.DebugLevel,
	config.LogLevelInfo:  applogger.InfoLevel,
	config.LogLevelWarn:  applogger.WarnLevel,
	config.LogLevelError: applogger.ErrorLevel,
	config.LogLevelFatal: applogger.FatalLevel,
	config.LogLevelPanic: applogger.PanicLevel,
}

func initLogger(cfg config.Config, closer *internalio.MultiCloser) (applogger.Logger, error) {
	loggerLevel, ok := configToLoggerLogLevelMap[cfg.LogLevel]
	if !ok {
		return nil, fmt.Errorf("unknown log level: %s", cfg.LogLevel)
	}

	logWriter, err := initLogWriter(cfg, closer)
	if err != nil {
		return nil, fmt.Errorf("error initializing log writer: %w", err)
	}

	logger, err := applogger.NewZapLogger(loggerLevel, logWriter)
	if err != nil {
		return nil, fmt.Errorf("error init logger: %w", err)
	}
	return logger, nil
}

func initLogWriter(cfg config.Config, closer *internalio.MultiCloser) (io.Writer, error) {
	if cfg.LogFile != nil && *cfg.LogFile != "" {
		file, err := os.OpenFile(*cfg.LogFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return nil, fmt.Errorf("error opening log file: %w", err)
		}

		if closer != nil {
			closer.AddCloser(internalio.CloserFunc(
				func() error {
					return file.Close()
				},
			))
		}

		return file, nil
	}

	return os.Stderr, nil
}

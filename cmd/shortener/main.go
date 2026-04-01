package main

import (
	"io"
	"log"
	"net/http"
	"os"

	internalio "github.com/liebeSonne/shortlink/internal/io"
	applogger "github.com/liebeSonne/shortlink/internal/logger"

	"github.com/liebeSonne/shortlink/internal/config"
	"github.com/liebeSonne/shortlink/internal/handler"
	"github.com/liebeSonne/shortlink/internal/handler/compress"
	"github.com/liebeSonne/shortlink/internal/repository"
	"github.com/liebeSonne/shortlink/internal/service"
)

const appID = "shortlink"
const envPrefix = ""

func main() {
	closer := internalio.MultiCloser{}
	defer func() {
		closeErr := closer.Close()
		if closeErr != nil {
			log.Fatalf("error closing closer: %v", closeErr)
		}
	}()

	cfg := initConfig()
	logger := initLogger(cfg, &closer)

	err := runApp(cfg, logger, &closer)

	logger.Fatalw("error starting server", "error", err)
}

func runApp(cfg config.Config, logger applogger.Logger, closer *internalio.MultiCloser) (err error) {
	shortLinkRepository := initShortLinkRepository(cfg, closer)
	shortIDGenerator := service.NewShortIDGenerator()
	shortLinkService := service.NewShortLinkService(shortLinkRepository, shortIDGenerator, service.DefaultMaxAttemptsToGenerateUniqueID)
	shortLinkHandler := handler.NewShortLinkHandler(shortLinkService, shortLinkRepository, cfg.BaseURL)
	databaseHandler := handler.NewDatabaseHandler()
	rootRouter := handler.NewRootRouter(shortLinkHandler, databaseHandler, cfg.EnableLogs)

	router := rootRouter.Router().(http.Handler)
	router, err = compress.NewCompressorMiddleware(router, compress.CompressorConfig{
		Encodings:    []compress.Encoding{compress.GzipEncoding},
		ContentTypes: &[]string{"application/json", "text/html"},
	})
	if err != nil {
		return err
	}
	router = handler.LoggingMiddleware(router, logger)

	logger.Infow("starting server",
		"addr", cfg.ServerAddress,
		"baseURL", cfg.BaseURL,
		"logLevel", cfg.LogLevel,
		"logFile", cfg.LogFile,
		"storage", cfg.FileStoragePath,
	)

	return http.ListenAndServe(cfg.ServerAddress, router)
}

var configToLoggerLogLevelMap = map[string]applogger.LogLevel{
	config.LogLevelDebug: applogger.DebugLevel,
	config.LogLevelInfo:  applogger.InfoLevel,
	config.LogLevelWarn:  applogger.WarnLevel,
	config.LogLevelError: applogger.ErrorLevel,
	config.LogLevelFatal: applogger.FatalLevel,
	config.LogLevelPanic: applogger.PanicLevel,
}

func initConfig() config.Config {
	cfg, err := config.LoadConfig(appID, envPrefix)
	if err != nil {
		log.Fatalf("error get config: %s", err.Error())
	}
	return cfg
}

func initLogger(cfg config.Config, closer *internalio.MultiCloser) applogger.Logger {
	loggerLevel, ok := configToLoggerLogLevelMap[cfg.LogLevel]
	if !ok {
		log.Fatalf("unknown log level: %s", cfg.LogLevel)
	}

	logWriter := initLogWriter(cfg, closer)

	logger, err := applogger.NewZapLogger(loggerLevel, logWriter)
	if err != nil {
		log.Fatalf("error init logger: %s", err.Error())
	}
	return logger
}

func initLogWriter(cfg config.Config, closer *internalio.MultiCloser) io.Writer {
	if cfg.LogFile != nil && *cfg.LogFile != "" {
		file, err := os.OpenFile(*cfg.LogFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}

		if closer != nil {
			closer.AddCloser(internalio.CloserFunc(
				func() error {
					return file.Close()
				},
			))
		}

		return file
	}

	return os.Stderr
}

func initShortLinkRepository(
	cfg config.Config,
	closer *internalio.MultiCloser,
) repository.ShortLinkRepository {
	if cfg.FileStoragePath != nil && *cfg.FileStoragePath != "" {
		repo, err := repository.NewFileShortLinkRepository(*cfg.FileStoragePath)
		if err != nil {
			log.Fatalf("error on init short link repository: %s", err.Error())
		}

		if closer != nil {
			closer.AddCloser(internalio.CloserFunc(
				func() error {
					return repo.Close()
				},
			))
		}

		return repo
	}

	return repository.NewMemoryShortLinkRepository()
}

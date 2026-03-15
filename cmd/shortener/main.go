package main

import (
	"log"
	"net/http"
	"os"

	"github.com/liebeSonne/shortlink/internal/config"
	"github.com/liebeSonne/shortlink/internal/handler"
	applogger "github.com/liebeSonne/shortlink/internal/logger"
	"github.com/liebeSonne/shortlink/internal/model"
	"github.com/liebeSonne/shortlink/internal/repository"
	"github.com/liebeSonne/shortlink/internal/service"
)

const appID = "shortlink"
const envPrefix = ""

func main() {
	cfg := initConfig()
	logger := initLogger(cfg)

	shortLinkRepository := repository.NewMemoryShortLinkRepository()
	shortIDGenerator := model.NewShortIDGenerator()
	shortLinkService := service.NewShortLinkService(shortLinkRepository, shortIDGenerator)
	shortLinkHandler := handler.NewShortLinkHandler(shortLinkService, shortLinkRepository, cfg.BaseURL)
	rootRouter := handler.NewRootRouter(shortLinkHandler, cfg.EnableLogs)
	router := handler.LoggingMiddleware(rootRouter.Router(), logger)

	logger.Infow("starting server", "addr", cfg.ServerAddress)
	err := http.ListenAndServe(cfg.ServerAddress, router)
	if err != nil {
		logger.Fatalw("error starting server", "error", err)
	}
}

var configToLoggerLogLevelMap = map[string]applogger.LogLevel{
	config.LogLevelDebug: applogger.InfoLevel,
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

func initLogger(cfg config.Config) applogger.Logger {
	loggerLevel, ok := configToLoggerLogLevelMap[cfg.LogLevel]
	if !ok {
		log.Fatalf("unknown log level: %s", cfg.LogLevel)
	}
	logger, err := applogger.NewZapLogger(loggerLevel, os.Stderr)
	if err != nil {
		log.Fatalf("error init logger: %s", err.Error())
	}
	return logger
}

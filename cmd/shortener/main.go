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
	conf, err := config.LoadConfig(appID, envPrefix)
	if err != nil {
		log.Fatalf("error get config: %s", err.Error())
	}

	loggerLevel := applogger.DebugLevel
	logger, err := applogger.NewZapLogger(loggerLevel, os.Stderr)
	if err != nil {
		log.Fatalf("error init logger: %s", err.Error())
	}

	shortLinkRepository := repository.NewMemoryShortLinkRepository()
	shortIDGenerator := model.NewShortIDGenerator()
	shortLinkService := service.NewShortLinkService(shortLinkRepository, shortIDGenerator)
	shortLinkHandler := handler.NewShortLinkHandler(shortLinkService, shortLinkRepository, conf.BaseURL)
	rootRouter := handler.NewRootRouter(shortLinkHandler, conf.EnableLogs)
	router := handler.LoggingMiddleware(rootRouter.Router(), logger)

	logger.Infow("starting server", "addr", conf.ServerAddress)
	err = http.ListenAndServe(conf.ServerAddress, router)
	if err != nil {
		logger.Fatalw("error starting server", "error", err)
	}
}

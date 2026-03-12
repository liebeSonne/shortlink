package main

import (
	"log"
	"net/http"

	"github.com/liebeSonne/shortlink/internal/config"
	"github.com/liebeSonne/shortlink/internal/handler"
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

	shortLinkRepository := repository.NewMemoryShortLinkRepository()
	shortIDGenerator := model.NewShortIDGenerator()
	shortLinkService := service.NewShortLinkService(shortLinkRepository, shortIDGenerator)
	shortLinkHandler := handler.NewShortLinkHandler(shortLinkService, shortLinkRepository, conf.BaseURL)
	rootRouter := handler.NewRootRouter(shortLinkHandler, conf.EnableLogs)

	err = http.ListenAndServe(conf.ServerAddress, rootRouter.Router())
	if err != nil {
		panic(err)
	}
}

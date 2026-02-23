package main

import (
	"net/http"

	"github.com/liebeSonne/shortlink/internal/handler"
	"github.com/liebeSonne/shortlink/internal/model"
	"github.com/liebeSonne/shortlink/internal/repository"
	"github.com/liebeSonne/shortlink/internal/service"
)

func main() {
	shortLinkRepository := repository.NewMemoryShortLinkRepository()
	shortIDGenerator := model.NewShortIDGenerator()
	shortLinkService := service.NewShortLinkService(shortLinkRepository, shortIDGenerator)
	shortLinkHandler := handler.NewShortLinkHandler(shortLinkService, shortLinkRepository)
	rootRouter := handler.NewRootRouter(shortLinkHandler)

	err := http.ListenAndServe(":8080", rootRouter.Router())
	if err != nil {
		panic(err)
	}
}

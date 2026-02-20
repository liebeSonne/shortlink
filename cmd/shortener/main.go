package main

import (
	"net/http"

	"github.com/liebeSonne/shortlink/internal/handler"
	"github.com/liebeSonne/shortlink/internal/model"
	"github.com/liebeSonne/shortlink/internal/repository"
	"github.com/liebeSonne/shortlink/internal/service"
)

type shortLink struct {
	id   string
	link string
}

func main() {
	shortLinkRepository := repository.NewMemoryShortLinkRepository()
	shortIDGenerator := model.NewShortIDGenerator()
	shortLinkService := service.NewShortLinkService(shortLinkRepository, shortIDGenerator)
	shortLinkHandler := handler.NewShortLinkHandler(shortLinkService, shortLinkRepository)
	rootHandler := handler.NewRootHandler(shortLinkHandler)

	mux := http.NewServeMux()
	mux.HandleFunc("/", rootHandler.Handle)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}

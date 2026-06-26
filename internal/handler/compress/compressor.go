package compress

import (
	"net/http"

	"github.com/liebeSonne/shortlink/internal/logger"
)

// CompressorConfig - настройки компрессоров.
type CompressorConfig struct {
	Encodings    []Encoding
	ContentTypes *[]string
}

// NewCompressorMiddleware - создание экземпляра посредника компрессоров.
func NewCompressorMiddleware(h http.Handler, cfg CompressorConfig, logger logger.Logger) (http.Handler, error) {
	for _, encoding := range cfg.Encodings {
		next, err := NewEncodingMiddleware(h, encoding, cfg.ContentTypes, logger)
		if err != nil {
			return nil, err
		}
		h = next
	}
	return h, nil
}

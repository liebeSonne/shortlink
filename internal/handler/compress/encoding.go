package compress

import (
	"errors"
	"net/http"

	"github.com/liebeSonne/shortlink/internal/handler/compress/deflate"
	"github.com/liebeSonne/shortlink/internal/handler/compress/gzip"
	"github.com/liebeSonne/shortlink/internal/logger"
)

var ErrUnknownEncoding = errors.New("unknown encoding")

// Encoding - тип кодировщика.
type Encoding int

// Значения кодировщиков.
const (
	GzipEncoding Encoding = iota
	DeflateEncoding
)

// NewEncodingMiddleware - создание экземпляра посредников кодировщиков.
func NewEncodingMiddleware(
	h http.Handler,
	encoding Encoding,
	contentTypes *[]string,
	logger logger.Logger,
) (http.HandlerFunc, error) {
	switch encoding {
	case GzipEncoding:
		return gzip.NewGzipHandlerMiddleware(h, contentTypes, logger), nil
	case DeflateEncoding:
		return deflate.NewDeflateHandlerMiddleware(h, contentTypes, logger), nil
	default:
		return nil, ErrUnknownEncoding
	}
}

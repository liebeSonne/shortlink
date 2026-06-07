package compress

import (
	"errors"
	"net/http"

	"github.com/liebeSonne/shortlink/internal/handler/compress/deflate"
	"github.com/liebeSonne/shortlink/internal/handler/compress/gzip"
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
) (http.HandlerFunc, error) {
	switch encoding {
	case GzipEncoding:
		return gzip.NewGzipHandlerMiddleware(h, contentTypes), nil
	case DeflateEncoding:
		return deflate.NewDeflateHandlerMiddleware(h, contentTypes), nil
	default:
		return nil, ErrUnknownEncoding
	}
}

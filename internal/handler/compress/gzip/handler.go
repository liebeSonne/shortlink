package gzip

import (
	"net/http"
	"slices"
	"strings"

	"github.com/liebeSonne/shortlink/internal/logger"
)

// NewGzipHandlerMiddleware - создание экземпляра посредника для gzip.
func NewGzipHandlerMiddleware(h http.Handler, contentTypes *[]string, logger logger.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		allowedContentType := true
		if contentTypes != nil {
			headerContentTypes := r.Header.Values("Content-Type")
			allowedContentType = slices.ContainsFunc(*contentTypes, func(contentType string) bool {
				return slices.ContainsFunc(headerContentTypes, func(s string) bool {
					return strings.Contains(s, contentType)
				})
			})
		}

		writer := w

		if allowedContentType {
			acceptGzip := slices.ContainsFunc(r.Header.Values("Accept-Encoding"), func(s string) bool {
				return strings.Contains(s, "gzip")
			})

			if acceptGzip {
				cw := NewGzipWriter(writer)
				defer func() {
					err := cw.Close()
					if err != nil {
						logger.Errorf("error closing gzip writer: %v", err)
					}
				}()
				writer = cw
			}
		}

		contentEncoding := slices.ContainsFunc(r.Header.Values("Content-Encoding"), func(s string) bool {
			return strings.Contains(s, "gzip")
		})

		if contentEncoding {
			cr, err := NewGzipReader(r.Body)
			if err != nil {
				logger.Errorf("error creating gzip reader: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			defer func() {
				err := cr.Close()
				if err != nil {
					logger.Errorf("error closing gzip writer: %v", err)
				}
			}()
			r.Body = cr
		}

		h.ServeHTTP(writer, r)
	}
}

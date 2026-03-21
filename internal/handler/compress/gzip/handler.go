package gzip

import (
	"fmt"
	"net/http"
	"slices"
	"strings"
)

func NewGzipHandlerMiddleware(h http.Handler, contentTypes *[]string) http.HandlerFunc {
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

		if !allowedContentType {
			h.ServeHTTP(w, r)
			return
		}

		writer := w

		acceptGzip := slices.ContainsFunc(r.Header.Values("Accept-Encoding"), func(s string) bool {
			return strings.Contains(s, "gzip")
		})

		if acceptGzip {
			cw := NewGzipWriter(writer)
			defer func() {
				err := cw.Close()
				if err != nil {
					fmt.Printf("error closing gzip writer: %v\n", err)
				}
			}()
			writer = cw
		}

		contentEncoding := slices.ContainsFunc(r.Header.Values("Content-Encoding"), func(s string) bool {
			return strings.Contains(s, "gzip")
		})

		if contentEncoding {
			cr, err := NewGzipReader(r.Body)
			if err != nil {
				fmt.Printf("error creating gzip reader: %v\n", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			defer func() {
				err := cr.Close()
				if err != nil {
					fmt.Printf("error closing gzip writer: %v\n", err)
				}
			}()
			r.Body = cr
		}

		h.ServeHTTP(writer, r)
	}
}

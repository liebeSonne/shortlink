package deflate

import (
	"fmt"
	"net/http"
	"slices"
	"strings"
)

func NewDeflateHandlerMiddleware(h http.Handler, contentTypes *[]string) http.HandlerFunc {
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
				return strings.Contains(s, "deflate")
			})

			if acceptGzip {
				cw, err := NewDeflateWriter(writer)
				if err != nil {
					fmt.Printf("error creating deflate writer: %v\n", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				defer func() {
					err := cw.Close()
					if err != nil {
						fmt.Printf("error closing deflate writer: %v\n", err)
					}
				}()
				writer = cw
			}
		}

		contentEncoding := slices.ContainsFunc(r.Header.Values("Content-Encoding"), func(s string) bool {
			return strings.Contains(s, "deflate")
		})

		if contentEncoding {
			cr, err := NewDeflateReader(r.Body)
			if err != nil {
				fmt.Printf("error creating deflate reader: %v\n", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			defer func() {
				err := cr.Close()
				if err != nil {
					fmt.Printf("error closing deflate writer: %v\n", err)
				}
			}()
			r.Body = cr
		}

		h.ServeHTTP(writer, r)
	}
}

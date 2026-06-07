package compress

import "net/http"

// CompressorConfig - настройки компрессоров.
type CompressorConfig struct {
	Encodings    []Encoding
	ContentTypes *[]string
}

// NewCompressorMiddleware - создание экземпляра посредника компрессоров.
func NewCompressorMiddleware(h http.Handler, cfg CompressorConfig) (http.Handler, error) {
	for _, encoding := range cfg.Encodings {
		next, err := NewEncodingMiddleware(h, encoding, cfg.ContentTypes)
		if err != nil {
			return nil, err
		}
		h = next
	}
	return h, nil
}

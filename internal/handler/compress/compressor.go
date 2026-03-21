package compress

import "net/http"

type CompressorConfig struct {
	Encodings    []Encoding
	ContentTypes *[]string
}

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

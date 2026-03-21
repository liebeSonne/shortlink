package gzip

import (
	"compress/gzip"
	"io"
	"net/http"
)

type Writer interface {
	http.ResponseWriter
	io.Closer
}

func NewGzipWriter(w http.ResponseWriter) Writer {
	return &gzipWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

type gzipWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

func (w *gzipWriter) Header() http.Header {
	return w.w.Header()
}

func (w *gzipWriter) Write(bytes []byte) (int, error) {
	return w.zw.Write(bytes)
}

func (w *gzipWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		w.w.Header().Set("Content-Encoding", "gzip")
	}
	w.w.WriteHeader(statusCode)
}

func (w *gzipWriter) Close() error {
	return w.zw.Close()
}

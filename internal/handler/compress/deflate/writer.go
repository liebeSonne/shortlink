package deflate

import (
	"compress/zlib"
	"io"
	"net/http"
)

type Writer interface {
	http.ResponseWriter
	io.Closer
}

func NewDeflateWriter(w http.ResponseWriter) (Writer, error) {
	zw, err := zlib.NewWriterLevel(w, zlib.BestSpeed)
	if err != nil {
		return nil, err
	}
	return &zlibWriter{
		w:  w,
		zw: zw,
	}, nil
}

type zlibWriter struct {
	w  http.ResponseWriter
	zw *zlib.Writer
}

func (w *zlibWriter) Header() http.Header {
	return w.w.Header()
}

func (w *zlibWriter) Write(bytes []byte) (int, error) {
	return w.zw.Write(bytes)
}

func (w *zlibWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		w.w.Header().Set("Content-Encoding", "deflate")
	}
	w.w.WriteHeader(statusCode)
}

func (w *zlibWriter) Close() error {
	return w.zw.Close()
}

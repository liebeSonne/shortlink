package handler

import (
	"net/http"
	"time"

	"github.com/liebeSonne/shortlink/internal/logger"
)

func LoggingMiddleware(next http.Handler, logger logger.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		respData := responseData{
			statusCode: 0,
			size:       0,
		}
		lw := loggingResponseWriter{
			next:         w,
			responseData: &respData,
		}

		next.ServeHTTP(&lw, r)

		duration := time.Since(start)

		logger.Infow("",
			"uri", r.RequestURI,
			"method", r.Method,
			"status", respData.statusCode,
			"duration", duration,
			"size", respData.size,
		)
	})
}

type responseData struct {
	statusCode int
	size       int
}

type loggingResponseWriter struct {
	next         http.ResponseWriter
	responseData *responseData
}

func (r *loggingResponseWriter) Header() http.Header {
	return r.next.Header()
}

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.next.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.next.WriteHeader(statusCode)
	r.responseData.statusCode = statusCode
}

package middleware

import (
	"context"
	"github.com/zelas91/gofermart/internal/logger"
	"github.com/zelas91/gofermart/internal/types"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func Logger(log *zap.SugaredLogger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), types.Logger, log))
			next.ServeHTTP(w, r)

		})
	}
}

type loggingWriter struct {
	http.ResponseWriter
	status int
	size   int
}

func (l *loggingWriter) Write(b []byte) (int, error) {
	size, err := l.ResponseWriter.Write(b)
	l.size = size
	return size, err
}

func (l *loggingWriter) WriteHeader(statusCode int) {
	l.status = statusCode
	l.ResponseWriter.WriteHeader(statusCode)
}

func WithLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := logger.GetLogger(r.Context())
		logWriter := &loggingWriter{
			ResponseWriter: w,
			status:         200,
		}

		start := time.Now()
		next.ServeHTTP(logWriter, r)
		duration := time.Since(start)

		log.Debugf("url : %s, method : %s, duration : %s status : %d, size : %d, content-type : %s",
			r.URL, r.Method, duration, logWriter.status, logWriter.size, r.Header.Get("Content-Type"))
	})

}

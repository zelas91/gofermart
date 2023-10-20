package middleware

import (
	"context"
	"github.com/zelas91/gofermart/internal/types"
	"go.uber.org/zap"
	"net/http"
)

func Logger(log *zap.SugaredLogger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), types.Logger, log))
			next.ServeHTTP(w, r)

		})
	}
}

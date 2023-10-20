package middleware

import (
	"context"
	"github.com/zelas91/gofermart/internal/logger"
	"github.com/zelas91/gofermart/internal/payload"
	"github.com/zelas91/gofermart/internal/service"
	"github.com/zelas91/gofermart/internal/types"
	"net/http"
)

func ValidationAuthorization(authService service.Authorization) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := logger.GetLogger(r.Context())

			cookie, err := r.Cookie("jwt")
			if err != nil {
				log.Errorf("not found jwt (err : %v)", err)
				payload.NewErrorResponse(w, "not found jwt", http.StatusUnauthorized)
				return
			}
			user, err := authService.ParserToken(r.Context(), cookie.Value)
			if err != nil {
				log.Errorf("parse token err : %v", err)
				payload.NewErrorResponse(w, err.Error(), http.StatusUnauthorized)
				return
			}
			r = r.WithContext(context.WithValue(r.Context(), types.UserIDKey, user.ID))
			next.ServeHTTP(w, r)
		})
	}
}

package middleware

import (
	"context"
	"github.com/zelas91/gofermart/internal/payload"
	"github.com/zelas91/gofermart/internal/service"
	"net/http"
)

func ValidationAuthorization(authService service.Authorization) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("jwt")
			if err != nil {
				payload.NewErrorResponse(w, "not found jwt", http.StatusUnauthorized)
				return
			}
			user, err := authService.ParserToken(r.Context(), cookie.Value)
			if err != nil {
				payload.NewErrorResponse(w, err.Error(), http.StatusUnauthorized)
				return
			}
			ctx := r.Context()
			ctx = context.WithValue(ctx, "userId", user.ID)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

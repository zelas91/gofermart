package controllers

import (
	"github.com/go-chi/chi/v5"
	"github.com/zelas91/gofermart/internal/logger"
	"github.com/zelas91/gofermart/internal/middleware"
	"github.com/zelas91/gofermart/internal/service"
	"github.com/zelas91/gofermart/internal/types"
	"go.uber.org/zap"
	"io"
	"net/http"
)

var content = "Content-Type"

type Handler struct {
	services *service.Service
	accrual  string
}

func NewHandler(accrual string, services *service.Service) *Handler {
	return &Handler{services: services, accrual: accrual}
}

func (h *Handler) InitRoutes(log *zap.SugaredLogger) http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Logger(log), middleware.WithLogging)
	router.Route("/api/user", func(r chi.Router) {
		r.Post("/register", h.signUp())
		r.Post("/login", h.signIn())
		r.Route("/", func(r chi.Router) {
			r.Use(middleware.ValidationAuthorization(h.services))
			r.Get("/orders", h.getOrders())
			r.Post("/orders", h.postOrders())
			r.Get("/balance", pass())
			r.Post("/balance/withdraw", pass())
			r.Get("/withdrawals", pass())
		})

	})
	return router
}

func pass() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		log := logger.GetLogger(request.Context())
		log.Info(request.Context().Value(types.UserIDKey))
		body, err := io.ReadAll(request.Body)

		log.Info(request.Method, " ", request.URL.Path, " ", string(body), " ", err)
	}
}

package controllers

import (
	"github.com/go-chi/chi/v5"
	middleware2 "github.com/go-chi/chi/v5/middleware"
	"github.com/zelas91/gofermart/internal/middleware"
	"github.com/zelas91/gofermart/internal/service"
	"go.uber.org/zap"
	"net/http"
)

var content = "Content-Type"

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes(log *zap.SugaredLogger) http.Handler {
	router := chi.NewRouter()
	router.Use(middleware2.Recoverer, middleware.Logger(log), middleware.WithLogging)
	router.Route("/api/user", func(r chi.Router) {
		r.Post("/register", h.signUp())
		r.Post("/login", h.signIn())
		r.Route("/", func(r chi.Router) {
			r.Use(middleware.ValidationAuthorization(h.services))
			r.Get("/orders", h.getOrders())
			r.Post("/orders", h.postOrders())
			r.Get("/balance", h.getBalance())
			r.Post("/balance/withdraw", h.withdraw())
			r.Get("/withdrawals", h.withdrawInfo())
		})

	})
	return router
}

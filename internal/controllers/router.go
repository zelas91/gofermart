package controllers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/zelas91/gofermart/internal/middleware"
	"github.com/zelas91/gofermart/internal/service"
	"net/http"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() http.Handler {
	router := chi.NewRouter()
	router.Route("/api/user", func(r chi.Router) {
		r.Post("/register", h.signUp())
		r.Post("/login", h.signIn())
		r.Route("/orders", func(r chi.Router) {
			r.Use(middleware.ValidationAuthorization(h.services))
			r.Get("/", pass())
			r.Post("/", pass())
		})
	})
	return router
}

func pass() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		fmt.Println("PASS")
		fmt.Println(request.Context().Value("userId"))
	}
}

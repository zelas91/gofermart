package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/zelas91/gofermart/internal/controllers"
	"net/http"
)

func InitRoutes() http.Handler {
	router := chi.NewRouter()
	router.Route("/api/user", func(r chi.Router) {
		r.Post("/register", controllers.SignUp())
		r.Post("/login", pass())
		r.Route("/orders", func(r chi.Router) {
			r.Get("/", pass())
			r.Post("/", pass())
		})
	})
	return router
}

func pass() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
	}
}

package controllers

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/zelas91/gofermart/internal/entities"
	"net/http"
)

func SignUp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		user := &entities.User{}
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			w.WriteHeader(400)
			return
		}

		validate := validator.New()
		if err := validate.Struct(user); err != nil {
			w.WriteHeader(501)
			return
		}

		w.WriteHeader(201)
	}
}

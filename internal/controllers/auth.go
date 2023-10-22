package controllers

import (
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/zelas91/gofermart/internal/entities"
	"github.com/zelas91/gofermart/internal/logger"
	"github.com/zelas91/gofermart/internal/payload"
	"github.com/zelas91/gofermart/internal/repository"
	"net/http"
)

func (h *Handler) signUp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.GetLogger(r.Context())

		if r.Header.Get(content) != "application/json" {
			log.Errorf("invalid content type")
			payload.NewErrorResponse(w, "invalid content type", http.StatusUnsupportedMediaType)
			return
		}

		user := &entities.User{}
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			log.Errorf("sigUp json decode err :%v", err)
			payload.NewErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}

		validate := validator.New()
		if err := validate.Struct(user); err != nil {
			log.Errorf("sigUp json validate err :%v", err)
			payload.NewErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := h.services.CreateUser(r.Context(), user); err != nil {
			if errors.Is(err, repository.ErrDuplicate) {
				log.Errorf("sigUp user duplicate err :%v", err)
				payload.NewErrorResponse(w, err.Error(), http.StatusConflict)
				return
			}
			log.Errorf("sigUp create user err :%v", err)
			payload.NewErrorResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}
		token, err := h.services.CreateToken(r.Context(), user)
		if err != nil {
			log.Errorf("sigUp create token err :%v", err)
			payload.NewErrorResponse(w, "invalid login or password", http.StatusUnauthorized)
			return
		}
		cookies := http.Cookie{
			Path:  "/",
			Name:  "jwt",
			Value: token,
		}
		http.SetCookie(w, &cookies)
		w.WriteHeader(http.StatusOK)

	}
}

func (h *Handler) signIn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.GetLogger(r.Context())

		if r.Header.Get(content) != "application/json" {
			log.Errorf("invalid content type")
			payload.NewErrorResponse(w, "invalid content type", http.StatusUnsupportedMediaType)
			return
		}

		user := &entities.User{}
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			log.Errorf("signIn json decode err:%v", err)
			payload.NewErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}

		validate := validator.New()
		if err := validate.Struct(user); err != nil {
			log.Errorf("signIn json validate err:%v", err)
			payload.NewErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}
		token, err := h.services.CreateToken(r.Context(), user)
		if err != nil {
			log.Errorf("signIn create token err:%v", err)
			payload.NewErrorResponse(w, "invalid login or password", http.StatusUnauthorized)
			return
		}
		cookies := http.Cookie{
			Path:  "/",
			Name:  "jwt",
			Value: token,
		}
		http.SetCookie(w, &cookies)
		w.WriteHeader(http.StatusOK)
	}
}

package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/zelas91/gofermart/internal/entities"
	"github.com/zelas91/gofermart/internal/logger"
	"github.com/zelas91/gofermart/internal/payload"
	"github.com/zelas91/gofermart/internal/repository"
	"net/http"
)

func (h *Handler) getBalance() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		balance, err := h.services.Balance.GetBalance(r.Context())
		if err != nil {
			logger.GetLogger(r.Context()).Errorf("get balance err : %v", err)
			payload.NewErrorResponse(w, "get balance err ", http.StatusInternalServerError)
			return
		}
		body, err := json.Marshal(&balance)
		if err != nil {
			logger.GetLogger(r.Context()).Errorf("json encode err : %v", err)
			payload.NewErrorResponse(w, "json encode error", http.StatusInternalServerError)
			return
		}
		if _, err = w.Write(body); err != nil {
			logger.GetLogger(r.Context()).Errorf("write body err : %v", err)
			payload.NewErrorResponse(w, "write body error", http.StatusInternalServerError)
			return
		}
	}
}

func (h *Handler) withdraw() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get(content) != "application/json" {
			logger.GetLogger(r.Context()).Errorf("invalid content type")
			payload.NewErrorResponse(w, "invalid content type", http.StatusUnsupportedMediaType)
			return
		}

		var withdraw entities.Withdraw
		if err := json.NewDecoder(r.Body).Decode(&withdraw); err != nil {
			logger.GetLogger(r.Context()).Errorf("json decode err : %v", err)
			payload.NewErrorResponse(w, "json decode err ", http.StatusInternalServerError)
			return
		}
		validate := validator.New()
		if err := validate.Struct(withdraw); err != nil {
			logger.GetLogger(r.Context()).Errorf("sigUp json validate err :%v", err)
			payload.NewErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}
		if !h.services.ValidateNumber(withdraw.Order) {
			logger.GetLogger(r.Context()).Errorf("invalid number order")
			payload.NewErrorResponse(w,
				fmt.Sprintf("invalid number order = %s", withdraw.Order), http.StatusUnprocessableEntity)
			return
		}
		if err := h.services.Balance.Withdraw(r.Context(), withdraw); err != nil {
			if errors.Is(err, repository.ErrNotEnoughFunds) {
				logger.GetLogger(r.Context()).Errorf(" err : %v number = %s", err, withdraw.Order)
				payload.NewErrorResponse(w, fmt.Sprintf(" err : %v number = %s", err, withdraw.Order),
					http.StatusPaymentRequired)
				return
			}
		}

	}
}
func (h *Handler) withdrawInfo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		withdraws, err := h.services.Balance.WithdrawInfo(r.Context())
		if err != nil {
			logger.GetLogger(r.Context()).Errorf("get withdraws  err: %v", err)
			payload.NewErrorResponse(w, "get withdraws err", http.StatusInternalServerError)
			return
		}
		if len(withdraws) == 0 {
			payload.NewErrorResponse(w, "no withdraws", http.StatusNoContent)
			return
		}

		body, err := json.Marshal(&withdraws)
		if err != nil {
			logger.GetLogger(r.Context()).Errorf("json encode err : %v", err)
			payload.NewErrorResponse(w, "json encode err", http.StatusInternalServerError)
			return
		}

		if _, err = w.Write(body); err != nil {
			logger.GetLogger(r.Context()).Errorf("write body err : %v", err)
			return
		}
	}
}

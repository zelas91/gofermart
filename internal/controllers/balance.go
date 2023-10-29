package controllers

import (
	"encoding/json"
	"github.com/zelas91/gofermart/internal/logger"
	"github.com/zelas91/gofermart/internal/payload"
	"net/http"
)

func (h *Handler) getBalance() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		balance, err := h.services.Balance.GetBalance(r.Context())
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

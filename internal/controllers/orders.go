package controllers

import (
	"fmt"
	"github.com/zelas91/gofermart/internal/logger"
	"github.com/zelas91/gofermart/internal/payload"
	"github.com/zelas91/gofermart/internal/types"
	"io"
	"net/http"
)

func (h *Handler) postOrders() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.GetLogger(r.Context())

		if r.Header.Get("Content-Type") != "text/plain" {
			log.Errorf("incorrect format request")
			payload.NewErrorResponse(w, "incorrect format request", http.StatusBadRequest)
			return
		}
		data, err := io.ReadAll(r.Body)
		if err != nil {
			log.Errorf("read body err : %v", err)
			payload.NewErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer func(body io.ReadCloser) {
			_ = body.Close()
		}(r.Body)
		number := string(data)

		if !h.services.ValidateNumber(number) {
			log.Errorf("invalid number order")
			payload.NewErrorResponse(w,
				fmt.Sprintf("invalid number order = %s", number), http.StatusUnprocessableEntity)
			return
		}

		authUserID := r.Context().Value(types.UserIDKey).(int64)

		userID, err := h.services.FindUserIDByOrder(r.Context(), number)
		if err != nil {
			log.Errorf("find user id by order err : %v", err)
			payload.NewErrorResponse(w, "find user id by order", http.StatusInternalServerError)
			return
		}
		if userID != 0 {
			if authUserID == userID {
				w.WriteHeader(http.StatusOK)
				return
			}
			payload.NewErrorResponse(w, "order will be loaded on another user", http.StatusConflict)
			return
		}

		if err = h.services.CreateOrder(r.Context(), number); err != nil {
			log.Errorf("create order err : %v", err)
			payload.NewErrorResponse(w, "create order err", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusAccepted)

	}
}

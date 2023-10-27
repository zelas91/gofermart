package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/zelas91/gofermart/internal/logger"
	"github.com/zelas91/gofermart/internal/payload"
	"io"
	"net/http"
)

func (h *Handler) postOrders() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Header.Get("Content-Type") != "text/plain" {
			logger.GetLogger(r.Context()).Errorf("incorrect format request")
			payload.NewErrorResponse(w, "incorrect format request", http.StatusBadRequest)
			return
		}
		data, err := io.ReadAll(r.Body)
		if err != nil {
			logger.GetLogger(r.Context()).Errorf("read body err : %v", err)
			payload.NewErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer func(body io.ReadCloser) {
			_ = body.Close()
		}(r.Body)
		number := string(data)

		if !h.services.ValidateNumber(number) {
			logger.GetLogger(r.Context()).Errorf("invalid number order")
			payload.NewErrorResponse(w,
				fmt.Sprintf("invalid number order = %s", number), http.StatusUnprocessableEntity)
			return
		}

		//authUserID := r.Context().Value(types.UserIDKey).(int64)
		//
		//userID, err := h.services.FindUserIDByOrder(r.Context(), number)
		//if err != nil {
		//	logger.GetLogger(r.Context()).Errorf("find user id by order err : %v", err)
		//	payload.NewErrorResponse(w, "find user id by order", http.StatusInternalServerError)
		//	return
		//}
		//if userID != 0 {
		//	if authUserID == userID {
		//		w.WriteHeader(http.StatusOK)
		//		return
		//	}
		//	payload.NewErrorResponse(w, "the order was uploaded to another user", http.StatusConflict)
		//	return
		//}
		if err = h.services.CreateOrder(r.Context(), number); err != nil {
			logger.GetLogger(r.Context()).Errorf("create order err : %v", err)
			payload.NewErrorResponse(w, "create order err", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusAccepted)
	}
}

func (h *Handler) getOrders() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		orders, err := h.services.FindOrdersByUser(r.Context())
		if len(orders) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		if err != nil {
			logger.GetLogger(r.Context()).Errorf("get orders err: %v", err)
			payload.NewErrorResponse(w, "get orders err", http.StatusInternalServerError)
			return
		}
		body, err := json.Marshal(&orders)
		if err != nil {
			logger.GetLogger(r.Context()).Errorf("orders encode to json : %v", err)
			payload.NewErrorResponse(w, "order encode to json err", http.StatusInternalServerError)
			return
		}

		if _, err = w.Write(body); err != nil {
			logger.GetLogger(r.Context()).Errorf("write body err : %v", err)
			payload.NewErrorResponse(w, "", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)

	}
}

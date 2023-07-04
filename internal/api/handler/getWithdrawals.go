package handler

import (
	"context"
	"encoding/json"
	"go-diploma/internal/api/models"
	"log"
	"net/http"
	"time"
)

func (h *LoyaltyHandler) GetWithdrawals(w http.ResponseWriter, r *http.Request) {
	if r.ContentLength > 0 {
		http.Error(w, "Content not expected", http.StatusBadRequest)
		return
	}

	uID := userID(r)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	withdrawals, errList := h.OrderSrv.ListUserWithdrawals(ctx, uID)
	if errList != nil {
		log.Printf("Cannot list withdrawals for user [%d], %s", uID, errList.Error())
		http.Error(w, "Cannot list withdrawals for user", http.StatusInternalServerError)
		return
	}

	view := make([]models.WithdrawalView, 0, len(withdrawals))
	for _, withdrawal := range withdrawals {
		view = append(view, models.NewWithdrawalView(withdrawal))
	}

	w.Header().Set("Content-Type", "application/json")
	if len(view) > 0 {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "    ")
	if err := enc.Encode(view); err != nil {
		log.Printf("Cannot write response: %s", err.Error())
	}
}

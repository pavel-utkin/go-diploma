package handler

import (
	"context"
	"encoding/json"
	"go-diploma/internal/api/models"
	"log"
	"net/http"
	"time"
)

func (h *LoyaltyHandler) Balance(w http.ResponseWriter, r *http.Request) {
	if r.ContentLength > 0 {
		http.Error(w, "Content not expected", http.StatusBadRequest)
		return
	}

	uID := userID(r)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	balance, errBalance := h.OrderSrv.GetUserBalance(uID, ctx)
	if errBalance != nil {
		log.Printf("Cannot get balance for user [%d]: %s", uID, errBalance.Error())
		http.Error(w, "Cannot get balance for user", http.StatusInternalServerError)
		return
	}

	view := models.NewBalanceView(balance)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	enc := json.NewEncoder(w)
	enc.SetIndent("", "    ")
	if err := enc.Encode(view); err != nil {
		log.Printf("Cannot write response: %s", err.Error())
	}
}

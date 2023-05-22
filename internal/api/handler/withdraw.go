package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go-diploma/internal/api/models"
	"go-diploma/internal/service/order"
	"log"
	"net/http"
	"time"
)

func (h *LoyaltyHandler) Withdraw(w http.ResponseWriter, r *http.Request) {
	if ct := r.Header.Get("Content-Type"); ct != "application/json" {
		http.Error(w, "Expected application/json content type", http.StatusBadRequest)
		return
	}

	j := models.WithdrawalRequestJSON{}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&j); err != nil {
		msg := fmt.Sprintf("Invalid withdrawal request json: %s", err.Error())
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	uID := userID(r)
	wr, errReq := models.NewWithdrawalRequest(j, uID)
	if errors.Is(errReq, order.ErrInvalidOrderNr) {
		msg := fmt.Sprintf("Invalid order nr: %s", errReq.Error())
		http.Error(w, msg, http.StatusUnprocessableEntity)
		return
	}
	if errReq != nil {
		log.Printf("Cannot parse withdrawal request: %s", errReq.Error())
		http.Error(w, "Cannot parse withdrawal request", http.StatusBadRequest)
		return
	}

	errWithdraw := h.OrderSrv.Withdraw(wr, ctx)
	if errors.Is(errWithdraw, order.ErrInsufficientBalance) {
		http.Error(w, "Insufficient balance", http.StatusPaymentRequired)
		return
	}
	if errWithdraw != nil {
		log.Printf("Cannot withdraw for user [%d]: %s", uID, errWithdraw.Error())
		http.Error(w, "Cannot withdraw", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

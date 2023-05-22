package handler

import (
	"context"
	"encoding/json"
	"go-diploma/internal/api/models"
	"log"
	"net/http"
	"time"
)

func (h *LoyaltyHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	if r.ContentLength > 0 {
		http.Error(w, "Content not expected", http.StatusBadRequest)
		return
	}

	uID := userID(r)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	orders, errList := h.OrderSrv.ListUserOrders(uID, ctx)
	if errList != nil {
		log.Printf("Cannot list orders for user [%d], %s", uID, errList.Error())
		http.Error(w, "Cannot list orders for user", http.StatusInternalServerError)
		return
	}

	view := make([]models.OrderView, 0, len(orders))
	for _, order := range orders {
		view = append(view, models.NewOrderView(order))
	}

	w.Header().Set("Content-Type", "application/json")
	if len(view) > 0 {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}

	log.Printf("Listing user orders: %v", view)

	enc := json.NewEncoder(w)
	enc.SetIndent("", "    ")
	if err := enc.Encode(view); err != nil {
		log.Printf("Cannot write response: %s", err.Error())
	}
}

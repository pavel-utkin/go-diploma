package handler

import (
	"context"
	"errors"
	"fmt"
	"go-diploma/internal/service/order"
	"io"
	"log"
	"net/http"
	"time"
)

func (h *LoyaltyHandler) PostOrder(w http.ResponseWriter, r *http.Request) {
	if contentType := r.Header.Get("Content-Type"); contentType != "text/plain" {
		msg := fmt.Sprintf("Unsupported content type [%s]", contentType)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	nrStr, errRead := io.ReadAll(r.Body)
	if errRead != nil {
		msg := fmt.Sprintf("Cannot read request body: %s", errRead.Error())
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	nr, errNr := order.ParseOrderNr(string(nrStr))
	if errNr != nil {
		msg := fmt.Sprintf("Invalid order nr: %s", errNr.Error())
		http.Error(w, msg, http.StatusUnprocessableEntity)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	uID := userID(r)

	o := order.ProcessRequest{
		Nr:     nr,
		UserID: uID,
	}

	errProcess := h.OrderSrv.UploadOrder(ctx, o)
	if errors.Is(errProcess, order.ErrDuplicateOrderForUser) {
		w.WriteHeader(http.StatusOK)
		return
	}
	if errors.Is(errProcess, order.ErrDuplicateOrderForAnotherUser) {
		http.Error(w, "order already posted by another user", http.StatusConflict)
		return
	}
	if errProcess != nil {
		msg := fmt.Sprintf("Cannot process order [%v]: %s", o, errProcess.Error())
		log.Println(msg)
		http.Error(w, "Cannot process order", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

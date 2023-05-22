package models

import (
	"fmt"
	"go-diploma/internal/service/order"
	"log"
	"strconv"
	"time"
)

type OrderView struct {
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	Accrual    float64   `json:"accrual"`
	UploadedAt time.Time `json:"uploaded_at"`
}

func NewOrderView(o order.Order) OrderView {
	return OrderView{
		Number:     fmt.Sprintf("%d", o.Nr),
		Status:     o.Status,
		Accrual:    float64(o.Accrual) / 100,
		UploadedAt: time.Time{},
	}
}

type BalanceView struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

func NewBalanceView(b order.Balance) BalanceView {
	return BalanceView{
		Current:   float64(b.Current) / 100,
		Withdrawn: float64(b.Withdrawn) / 100,
	}
}

type WithdrawalRequestJSON struct {
	OrderNr string  `json:"order"`
	Sum     float64 `json:"sum"`
}

func NewWithdrawalRequest(j WithdrawalRequestJSON, userID int64) (order.WithdrawalRequest, error) {
	wr := order.WithdrawalRequest{}

	orderNr, err := order.ParseOrderNr(j.OrderNr)
	if err != nil {
		return wr, fmt.Errorf("cannot make withdrawal request: %s", err)
	}

	wr.OrderNr = orderNr
	wr.Sum = int64(j.Sum * 100)
	wr.UserID = userID

	log.Printf("Parsed WithdrawalRequest: %v", wr)
	return wr, nil
}

type WithdrawalView struct {
	Order       string    `json:"order"`
	Sum         float64   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}

func NewWithdrawalView(w order.Withdrawal) WithdrawalView {
	return WithdrawalView{
		Order:       strconv.FormatInt(w.OrderNr, 10),
		Sum:         float64(w.Sum) / 100,
		ProcessedAt: w.RequestedAt,
	}
}

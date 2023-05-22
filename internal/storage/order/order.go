package order

import (
	"context"
	"go-diploma/internal/service/order"
)

type Storage interface {
	AddOrder(pr order.ProcessRequest, ctx context.Context) error
	GetOrderByNr(nr int64, ctx context.Context) (order.Order, error)
	ListUserOrders(userID int64, ctx context.Context) ([]order.Order, error)
	Withdraw(wr order.WithdrawalRequest, ctx context.Context) error
	ListUserWithdrawals(userID int64, ctx context.Context) ([]order.Withdrawal, error)
}

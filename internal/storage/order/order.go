package order

import (
	"context"
	"go-diploma/internal/service/order"
)

type Storage interface {
	AddOrder(ctx context.Context, pr order.ProcessRequest) error
	GetOrderByNr(ctx context.Context, nr int64) (order.Order, error)
	ListUserOrders(ctx context.Context, userID int64) ([]order.Order, error)
	Withdraw(ctx context.Context, wr order.WithdrawalRequest) error
	ListUserWithdrawals(ctx context.Context, userID int64) ([]order.Withdrawal, error)
}

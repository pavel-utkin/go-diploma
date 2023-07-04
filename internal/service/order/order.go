package order

import "context"

type Service interface {
	UploadOrder(ctx context.Context, pr ProcessRequest) error
	ListUserOrders(ctx context.Context, userID int64) ([]Order, error)
	GetUserBalance(ctx context.Context, userID int64) (Balance, error)
	Withdraw(ctx context.Context, wr WithdrawalRequest) error
	ListUserWithdrawals(ctx context.Context, userID int64) ([]Withdrawal, error)
}

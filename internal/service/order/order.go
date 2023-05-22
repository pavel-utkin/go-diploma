package order

import "context"

type Service interface {
	UploadOrder(pr ProcessRequest, ctx context.Context) error
	ListUserOrders(userID int64, ctx context.Context) ([]Order, error)
	GetUserBalance(userID int64, ctx context.Context) (Balance, error)
	Withdraw(wr WithdrawalRequest, ctx context.Context) error
	ListUserWithdrawals(userID int64, ctx context.Context) ([]Withdrawal, error)
}

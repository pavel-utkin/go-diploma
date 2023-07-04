package v1

import (
	"context"
	"errors"
	"fmt"
	srv "go-diploma/internal/service/order"
	storage "go-diploma/internal/storage/order"
	"log"
)

type Service struct {
	storage storage.Storage
}

func NewService(storage storage.Storage) (*Service, error) {
	if storage == nil {
		return nil, errors.New("storage required")
	}

	return &Service{storage}, nil
}

func (s *Service) UploadOrder(ctx context.Context, pr srv.ProcessRequest) error {

	if errAdd := s.storage.AddOrder(ctx, pr); errAdd != nil {
		if errors.Is(errAdd, srv.ErrDuplicateOrder) {
			dupO, errGet := s.storage.GetOrderByNr(ctx, pr.Nr)
			if errGet != nil {
				return fmt.Errorf("cannot get details of a duplicate order: %w", errGet)
			}

			if dupO.UserID == pr.UserID {
				return srv.ErrDuplicateOrderForUser
			} else {
				return srv.ErrDuplicateOrderForAnotherUser
			}
		}

		return fmt.Errorf("cannot schedule order for processing: %w", errAdd)
	}

	return nil
}

func (s *Service) ListUserOrders(ctx context.Context, userID int64) ([]srv.Order, error) {
	orders, err := s.storage.ListUserOrders(ctx, userID)
	if err != nil {
		return orders, fmt.Errorf("cannot list orders for user [%d]: %w", userID, err)
	}

	return orders, nil
}

func (s *Service) GetUserBalance(ctx context.Context, userID int64) (srv.Balance, error) {
	balance := srv.NewBalance()

	accruals, errAccruals := s.storage.ListUserOrders(ctx, userID)
	if errAccruals != nil {
		return balance, fmt.Errorf("cannot list accruals for user [%d]: %w", userID, errAccruals)
	}

	for _, accrual := range accruals {
		if accrual.Status == srv.StatusProcessed {
			balance.Current += accrual.Accrual
			balance.LatestAccrual = accrual.Nr
		}
	}

	withdrawals, errWithdrawals := s.storage.ListUserWithdrawals(ctx, userID)
	if errWithdrawals != nil {
		return balance, fmt.Errorf("cannot list withdrawals for user [%d]: %w", userID, errAccruals)
	}
	for _, withdrawal := range withdrawals {
		balance.Current -= withdrawal.Sum
		balance.Withdrawn += withdrawal.Sum
		balance.LatestWithdrawal = withdrawal.OrderNr
	}

	log.Printf("Balance calculated: %v", balance)
	return balance, nil
}

func (s *Service) Withdraw(ctx context.Context, wr srv.WithdrawalRequest) error {
	bal, errBal := s.GetUserBalance(ctx, wr.UserID)
	if errBal != nil {
		return fmt.Errorf("cannot get user balance [%v]: %w", wr, errBal)
	}

	if bal.Current < wr.Sum {
		return fmt.Errorf("insufficient balance [%v]: %w", wr, srv.ErrInsufficientBalance)
	}

	wr.LatestAccrual = bal.LatestAccrual
	wr.LatestWithdrawal = bal.LatestWithdrawal

	if err := s.storage.Withdraw(ctx, wr); err != nil {
		return fmt.Errorf("cannot withdraw [%v]: %w", wr, err)
	}

	return nil
}

func (s *Service) ListUserWithdrawals(ctx context.Context, userID int64) ([]srv.Withdrawal, error) {
	withdrawals, err := s.storage.ListUserWithdrawals(ctx, userID)
	if err != nil {
		return withdrawals, fmt.Errorf("cannot list withdrawals for user [%d]: %w", userID, err)
	}

	return withdrawals, nil
}

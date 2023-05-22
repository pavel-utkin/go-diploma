package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	srv "go-diploma/internal/service/order"
	"log"
)

type OrderStorage struct {
	*sql.DB
}

type Scannable interface {
	Scan(dest ...interface{}) error
}

func NewOrderStorage(db *sql.DB) (*OrderStorage, error) {
	if db == nil {
		return nil, errors.New("db should not be nil")
	}
	return &OrderStorage{db}, nil
}

func (s *OrderStorage) AddOrder(pr srv.ProcessRequest, ctx context.Context) error {
	row := s.QueryRowContext(ctx, `
		insert into ORDERS (ORDERS_NR, USERS_ID, ORDERS_STATUS) 
		values($1, $2, $3) 
		returning ORDERS_NR, USERS_ID, ORDERS_STATUS, ORDERS_ACCRUAL, ORDERS_UPLOADED_AT
		`, pr.Nr, pr.UserID, srv.StatusNew)

	order := srv.Order{}

	err := mapOrder(&order, row)
	var dbErr *pgconn.PgError
	if errors.As(err, &dbErr) && dbErr.Code == pgerrcode.UniqueViolation {
		log.Printf("Duplicate order [%d]", pr.Nr)
		err = srv.ErrDuplicateOrder
	}
	if err != nil {
		return fmt.Errorf("cannot insert order: %w", err)
	}

	return nil
}

func (s *OrderStorage) GetOrderByNr(nr int64, ctx context.Context) (srv.Order, error) {
	row := s.QueryRowContext(ctx, `
		select ORDERS_NR, USERS_ID, ORDERS_STATUS, ORDERS_ACCRUAL, ORDERS_UPLOADED_AT
		from ORDERS
		where ORDERS_NR = $1
		`, nr)
	order := srv.Order{}

	if err := mapOrder(&order, row); err != nil {
		return order, fmt.Errorf("cannot select order: %w", err)
	}

	return order, nil
}

func (s *OrderStorage) ListUserWithdrawals(userID int64, ctx context.Context) ([]srv.Withdrawal, error) {
	result := make([]srv.Withdrawal, 0)

	rows, err := s.QueryContext(ctx, `
		select WITHDRAWALS_NR, USERS_ID, WITHDRAWALS_SUM, WITHDRAWALS_REQUESTED_AT
		from WITHDRAWALS
		where USERS_ID = $1
		order by WITHDRAWALS_REQUESTED_AT
	`,
		userID)
	if err != nil {
		return result, fmt.Errorf("cannot select withdrawals for user [%d]: %w", userID, err)
	}
	defer func(rows *sql.Rows) {
		if err := rows.Close(); err != nil {
			log.Printf("Cannot close result set: %s", err.Error())
		}
	}(rows)

	for rows.Next() {
		w := srv.Withdrawal{}
		if err := rows.Scan(&w.OrderNr, &w.UserID, &w.Sum, &w.RequestedAt); err != nil {
			return result, fmt.Errorf("cannot map all withdrawals from DB: %w", err)
		}

		result = append(result, w)
	}
	if rows.Err() != nil {
		return result, fmt.Errorf("cannot iterate all results from DB: %w", rows.Err())
	}

	return result, nil

}

func (s *OrderStorage) ListUserOrders(userID int64, ctx context.Context) ([]srv.Order, error) {
	result := make([]srv.Order, 0)

	rows, err := s.QueryContext(ctx, `
		select ORDERS_NR, USERS_ID, ORDERS_STATUS, ORDERS_ACCRUAL, ORDERS_UPLOADED_AT
		from ORDERS
		where USERS_ID = $1
		order by ORDERS_UPLOADED_AT
	`,
		userID)
	if err != nil {
		return result, fmt.Errorf("cannot select orders for user [%d]: %w", userID, err)
	}
	defer func(rows *sql.Rows) {
		if err := rows.Close(); err != nil {
			log.Printf("Cannot close result set: %s", err.Error())
		}
	}(rows)

	for rows.Next() {
		order := srv.Order{}

		if err := mapOrder(&order, rows); err != nil {
			return result, fmt.Errorf("cannot map all orders from DB: %w", err)
		}

		result = append(result, order)
	}
	if rows.Err() != nil {
		return result, fmt.Errorf("cannot iterate all results from DB: %w", rows.Err())
	}

	return result, nil
}

func (s *OrderStorage) Withdraw(wr srv.WithdrawalRequest, ctx context.Context) error {
	log.Printf("Processing withdrawal request: %v", wr)
	result, errExec := s.ExecContext(ctx, `
			with NEW_WITHDRAWAL as (
				select
					$1::bigint as WITHDRAWALS_NR,
					$2::bigint as USERS_ID,
					$3::bigint as WITHDRAWALS_SUM,
					coalesce((select ORDERS_NR from ORDERS where USERS_ID = $2 order by ORDERS_UPLOADED_AT desc limit 1), -1) as LATEST_ACCRUAL,
					coalesce((select WITHDRAWALS_NR from WITHDRAWALS where USERS_ID = $2 order by WITHDRAWALS_REQUESTED_AT desc limit 1), -1) as LATEST_WITHDRAWAL
			)
			insert into WITHDRAWALS
			select WITHDRAWALS_NR, USERS_ID, WITHDRAWALS_SUM
			from NEW_WITHDRAWAL
			where
					LATEST_ACCRUAL = $4
			  and LATEST_WITHDRAWAL = $5
			returning WITHDRAWALS_NR, USERS_ID, WITHDRAWALS_SUM
		`, wr.OrderNr, wr.UserID, wr.Sum, wr.LatestAccrual, wr.LatestWithdrawal)

	if errExec != nil {
		return fmt.Errorf("cannot insert withdrawal: %w", errExec)
	}

	affected, errAffected := result.RowsAffected()
	if errAffected != nil {
		return fmt.Errorf("cannot get affected rows: %w", errAffected)
	}
	if affected != 1 {
		return fmt.Errorf("withdrowal not accepted because of conflict")
	}

	return nil
}

func mapOrder(o *srv.Order, row Scannable) error {
	errScan := row.Scan(&o.Nr, &o.UserID, &o.Status, &o.Accrual, &o.UploadedAt)
	if errScan == sql.ErrNoRows {
		return srv.ErrOrderNotFound
	}
	if errScan != nil {
		return fmt.Errorf("cannot scan order from DB results: %w", errScan)
	}

	return nil
}

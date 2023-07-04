package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"go-diploma/internal/api/models"
	storage "go-diploma/internal/storage/accrual"
	"log"
)

type AccrualStorage struct {
	*sql.DB
}

func NewAccrualStorage(db *sql.DB) (*AccrualStorage, error) {
	if db == nil {
		return nil, errors.New("db should not be nil")
	}
	return &AccrualStorage{db}, nil
}

func (s *AccrualStorage) NextOrder(ctx context.Context) (int64, error) {
	row := s.QueryRowContext(ctx, `
		with NEXT_ORDERS as (
		    select ORDERS_NR
			from ORDERS
			where ORDERS_STATUS in ('NEW', 'PROCESSING')
				and ORDERS_PROCESSED_AT + interval '1 seconds' < current_timestamp
			order by ORDERS_PROCESSED_AT
		)
		update ORDERS
			set 
			    ORDERS_PROCESSED_AT = current_timestamp,
				ORDERS_STATUS = 'PROCESSING'
			where ORDERS_NR = (select ORDERS_NR from NEXT_ORDERS limit 1)
			returning ORDERS_NR
	`)

	var orderID int64

	err := row.Scan(&orderID)
	if errors.Is(err, sql.ErrNoRows) {
		return orderID, storage.ErrNoOrders
	}
	if err != nil {
		return orderID, fmt.Errorf("cannot map order ID: %w", err)
	}

	return orderID, nil
}

func (s *AccrualStorage) ApplyAccrual(o models.OrderAccrual, ctx context.Context) error {
	log.Printf("Applying order accrual: %v", o)
	result, errExec := s.ExecContext(ctx, `
		update ORDERS
		set ORDERS_ACCRUAL = $1,
		    ORDERS_STATUS = $2
		where ORDERS_NR = $3
			and ORDERS_STATUS = 'PROCESSING'
	`, o.Accrual, o.Status, o.OrderNr)
	if errExec != nil {
		return fmt.Errorf("cannot update order: %w", errExec)
	}

	affected, errAffected := result.RowsAffected()
	if errAffected != nil {
		return fmt.Errorf("cannot get affected rows: %w", errExec)
	}
	if affected != 1 {
		return fmt.Errorf("order not updated; expected 1 row affected, got %d", affected)
	}

	return nil
}

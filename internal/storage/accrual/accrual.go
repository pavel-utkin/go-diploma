package accrual

import (
	"context"
	"errors"
	"go-diploma/internal/api/models"
)

var ErrNoOrders = errors.New("no more orders to process")

type Storage interface {
	NextOrder(ctx context.Context) (int64, error)
	ApplyAccrual(o models.OrderAccrual, ctx context.Context) error
}

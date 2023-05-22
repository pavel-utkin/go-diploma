package models

import (
	"fmt"
	"time"
)

const StatusRegistered = "REGISTERED"
const StatusInvalid = "INVALID"
const StatusProcessing = "PROCESSING"
const StatusProcessed = "PROCESSED"

type OrderAccrual struct {
	OrderNr int64
	Status  string
	Accrual int64
}

type ErrTooManyRequests struct {
	RetryAfter time.Duration
	Err        error
}

func (e *ErrTooManyRequests) Error() string {
	wrapped := fmt.Errorf("too many requests; retry after %v sec: %w", e.RetryAfter.Seconds(), e.Err)
	return wrapped.Error()
}

func (e *ErrTooManyRequests) Unwrap() error {
	return e.Err
}

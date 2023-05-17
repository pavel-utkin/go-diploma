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

type errTooManyRequests struct {
	RetryAfter time.Duration
	Err        error
}

func (e *errTooManyRequests) Error() error {
	return fmt.Errorf("too many requests; retry after %v sec: %w", e.RetryAfter.Seconds(), e.Err)
}

func (e *errTooManyRequests) Unwrap() error {
	return e.Err
}

package order

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

var ErrDuplicateOrder = errors.New("duplicate order")

var ErrDuplicateOrderForUser = errors.New("order already posted by the user")

var ErrDuplicateOrderForAnotherUser = errors.New("order already posted by another user")

var ErrInvalidOrderNr = errors.New("invalid order number")

var ErrOrderNotFound = errors.New("order not found")

var ErrInsufficientBalance = errors.New("insufficient balance")

const StatusNew = "NEW"
const StatusProcessing = "PROCESSING"
const StatusInvalid = "INVALID"
const StatusProcessed = "PROCESSED"

func ParseOrderNr(s string) (int64, error) {
	nr, errConv := strconv.ParseInt(s, 10, 64)
	if errConv != nil {
		return nr, fmt.Errorf("cannot parse order number: %w", errConv)
	}

	n := len(s)
	checksum := 0

	for i := 1; i <= len(s); i++ {
		d, err := strconv.Atoi(string(s[n-i]))
		if err != nil {
			return nr, fmt.Errorf("contains non-digit character: %w", ErrInvalidOrderNr)
		}

		if i%2 == 0 {
			s := 2 * d
			if s > 9 {
				s -= 9
			}
			checksum += s
		} else {
			checksum += d
		}
	}

	if checksum%10 != 0 {
		return nr, fmt.Errorf("incorrect checksum: %w", ErrInvalidOrderNr)
	}

	return nr, nil
}

type ProcessRequest struct {
	Nr     int64
	UserID int64
}

type Order struct {
	UserID     int64
	Nr         int64
	Status     string
	Accrual    int64
	UploadedAt time.Time
}

type Balance struct {
	Current          int64
	Withdrawn        int64
	LatestAccrual    int64
	LatestWithdrawal int64
}

func NewBalance() Balance {
	return Balance{
		Current:          0,
		Withdrawn:        0,
		LatestAccrual:    -1,
		LatestWithdrawal: -1,
	}
}

type WithdrawalRequest struct {
	OrderNr          int64
	UserID           int64
	Sum              int64
	LatestAccrual    int64
	LatestWithdrawal int64
}

type Withdrawal struct {
	OrderNr     int64
	UserID      int64
	Sum         int64
	RequestedAt time.Time
}

package accrual

import (
	"fmt"
	"go-diploma/internal/api/models"
	"go-diploma/internal/service/order"
)

type OrderAccrualJSON struct {
	OrderNr string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}

func (j OrderAccrualJSON) ToOrderAccrual() (models.OrderAccrual, error) {
	nr, errNr := order.ParseOrderNr(j.OrderNr)
	if errNr != nil {
		return models.OrderAccrual{}, fmt.Errorf("envalid order nr [%s]: %w", j.OrderNr, errNr)
	}
	result := models.OrderAccrual{
		OrderNr: nr,
		Status:  j.Status,
		Accrual: int64(j.Accrual * 100),
	}

	if result.Status == models.StatusRegistered {
		result.Status = order.StatusProcessing
	} else if result.Status == models.StatusProcessing {
		result.Status = order.StatusProcessing
	} else if result.Status == models.StatusProcessed {
		result.Status = order.StatusProcessed
	} else {
		result.Status = order.StatusInvalid
	}

	return result, nil
}

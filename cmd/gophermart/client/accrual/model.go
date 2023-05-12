package accrual

type OrderAccrualJSON struct {
	OrderNr string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}

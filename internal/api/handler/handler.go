package handler

import "go-diploma/internal/service/auth"

type LoyaltyHandler struct {
	AuthSrv auth.Service
}

func NewHandler(authSrv auth.Service) (*LoyaltyHandler, error) {
	handler := LoyaltyHandler{
		AuthSrv: authSrv,
	}

	return &handler, nil
}

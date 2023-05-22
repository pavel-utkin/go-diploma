package api

import (
	"context"
	"fmt"
	"go-diploma/internal/api/handler"
	accrualClient "go-diploma/internal/client/accrual"
	"go-diploma/internal/service/accrual"
	auth "go-diploma/internal/service/auth/v1"
	order "go-diploma/internal/service/order/v1"
	accrualStorage "go-diploma/internal/storage/accrual"
	authStorage "go-diploma/internal/storage/auth"
	orderStorage "go-diploma/internal/storage/order"
	"log"
	"net/http"
)

type LoyaltyServer struct {
	http.Server
	accrualSrv *accrual.Service
}

func NewServer(
	addr string,
	accrualAddr string,
	authStorage authStorage.Storage,
	orderStorage orderStorage.Storage,
	accrualStorage accrualStorage.Storage,
) (*LoyaltyServer, error) {
	authSrv, errAuth := auth.NewService(authStorage)
	if errAuth != nil {
		return nil, fmt.Errorf("cannot get instance of Auth Service: %w", errAuth)
	}

	orderSrv, errOrder := order.NewService(orderStorage)
	if errOrder != nil {
		return nil, fmt.Errorf("cannot get instance of Order Service: %w", errOrder)
	}

	ctx := context.Background()

	acrClient := accrualClient.NewClient(accrualAddr)
	acrSrv, errAcr := accrual.NewService(acrClient, accrualStorage, ctx)
	if errAcr != nil {
		return nil, fmt.Errorf("cannot get instance of Accrual Service: %s", errAcr)
	}

	h, errHandler := handler.NewHandler(authSrv, orderSrv)
	if errHandler != nil {
		return nil, fmt.Errorf("cannot get instance of Handler: %w", errHandler)
	}

	r := newRouter(h)
	server := LoyaltyServer{
		Server: http.Server{
			Addr:    addr,
			Handler: r,
		},
		accrualSrv: acrSrv,
	}

	log.Println("Starting server...")

	go func() {
		err := server.ListenAndServe()
		if err != http.ErrServerClosed {
			log.Printf("Server failed: %s", err.Error())
		}
	}()

	return &server, nil
}

func (s *LoyaltyServer) Shutdown(ctx context.Context) error {
	if err := s.Server.Shutdown(ctx); err != nil {
		return err
	}

	log.Println("Server stopped.")

	return nil
}

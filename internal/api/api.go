package api

import (
	"context"
	"fmt"
	"go-diploma/internal/api/handler"
	auth "go-diploma/internal/service/auth/v1"
	authStorage "go-diploma/internal/storage/auth"
	"log"
	"net/http"
)

type LoyaltyServer struct {
	http.Server
}

func NewServer(
	addr string,
	accrualAddr string,
	authStorage authStorage.Storage,
) (*LoyaltyServer, error) {

	authSrv, errAuth := auth.NewService(authStorage)
	if errAuth != nil {
		return nil, fmt.Errorf("cannot get instance of Auth Service: %w", errAuth)
	}

	h, errHandler := handler.NewHandler(authSrv)
	if errHandler != nil {
		return nil, fmt.Errorf("cannot get instance of Handler: %w", errHandler)
	}

	r := newRouter(h)

	server := LoyaltyServer{
		Server: http.Server{
			Addr:    addr,
			Handler: r,
		},
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

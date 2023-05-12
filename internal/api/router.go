package api

import (
	"github.com/go-chi/chi/v5"
	"go-diploma/internal/api/handler"
)

type loyaltyRouter struct {
	*chi.Mux
}

func newRouter(h *handler.LoyaltyHandler) *loyaltyRouter {
	router := loyaltyRouter{
		Mux: chi.NewMux(),
	}

	router.Route("/api/user", func(r chi.Router) {
		r.Post("/register", h.Register)
		r.Post("/login", h.Login)
	})

	return &router
}

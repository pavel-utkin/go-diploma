package middleware

import (
	"context"
	"fmt"
	"go-diploma/internal/api/models"
	"go-diploma/internal/service/auth"
	"log"
	"net/http"
	"time"
)

type AuthContextKeyType struct{}

// Authenticator middleware authenticates a request
// based on the signed cookie containing a user ID.
// In case authentication has failed, it signs up a new user.
func Authenticator(s auth.Service) func(http.Handler) http.Handler {
	ra := requestAuth{s}

	return func(next http.Handler) http.Handler {
		serveHTTP := func(w http.ResponseWriter, r *http.Request) {
			userID := ra.extractUserID(r)
			if userID == nil {
				http.Error(w, "Login to access this endpoint", http.StatusUnauthorized)
				return
			}

			ctxWithUserID := context.WithValue(r.Context(), AuthContextKeyType{}, *userID)

			next.ServeHTTP(w, r.WithContext(ctxWithUserID))
		}

		return http.HandlerFunc(serveHTTP)
	}
}

type requestAuth struct {
	AuthService auth.Service
}

func (a *requestAuth) extractUserID(r *http.Request) *int64 {
	cookie, errGetCookie := r.Cookie(models.AuthCookieName)
	if errGetCookie != nil {
		log.Printf("Cannot get authentication cookie: %s", errGetCookie.Error())
		return nil
	}

	var userID int64
	var signature []byte
	if _, err := fmt.Sscanf(cookie.Value, "%d|%x", &userID, &signature); err != nil {
		log.Printf("Cannot parse authentication cookie [%s]: %s", cookie.Value, err.Error())
	}

	sgn := auth.SignedUserID{
		ID:        userID,
		Signature: signature,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if invalid := a.AuthService.Validate(sgn, ctx); invalid != nil {
		log.Printf("Signature is invalid: %s", invalid.Error())
		return nil
	}

	return &sgn.ID
}

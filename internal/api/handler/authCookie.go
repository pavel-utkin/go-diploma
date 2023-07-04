package handler

import (
	"fmt"
	"go-diploma/internal/api/models"
	"go-diploma/internal/service/auth"
	"log"
	"net/http"
)

func makeAuthCookie(u auth.SignedUserID) http.Cookie {
	v := fmt.Sprintf("%d|%x", u.ID, u.Signature)
	log.Println("Cookie : ", v)
	return http.Cookie{
		Name:  models.AuthCookieName,
		Value: v,
	}
}

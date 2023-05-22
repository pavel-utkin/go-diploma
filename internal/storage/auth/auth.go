package auth

import (
	"context"
	srv "go-diploma/internal/service/auth"
)

type Storage interface {
	CreateUser(u srv.UserToCreate, ctx context.Context) (srv.User, error)
	GetUserByLogin(login string, ctx context.Context) (*srv.User, error)
	SetUserSession(u srv.UserSessionToStart, ctx context.Context) (srv.UserSession, error)
	GetUserSession(uID int64, ctx context.Context) (srv.UserSession, error)
}

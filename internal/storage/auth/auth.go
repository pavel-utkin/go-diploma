package auth

import srv "go-diploma/internal/service/auth"

type Storage interface {
	CreateUser(u srv.UserToCreate) (srv.User, error)
	GetUserByLogin(login string) (*srv.User, error)
	SetUserSession(u srv.UserSessionToStart) (srv.UserSession, error)
	GetUserSession(uID int64) (srv.UserSession, error)
}

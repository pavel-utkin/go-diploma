package auth

import (
	"errors"
	"time"
)

var ErrLoginAlreadyTaken = errors.New("login already taken")

var ErrWrongCredentials = errors.New("incorrect login/password")

var ErrUserNotFound = errors.New("user not found")

type Credentials struct {
	Login    string
	Password []byte
}

type SignedUserID struct {
	ID        int64
	Signature []byte
}

type UserToCreate struct {
	Login        string
	PasswordHash []byte
}

type User struct {
	ID           int64
	Login        string
	PasswordHash []byte
}

type UserSessionToStart struct {
	UserID       int64
	SignatureKey []byte
}

type UserSession struct {
	UserID       int64
	SignatureKey []byte
	StartedAt    time.Time
}

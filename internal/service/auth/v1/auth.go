package v1

import (
	"context"
	"crypto/hmac"
	"errors"
	"fmt"
	"go-diploma/internal/service/auth"
	authStorage "go-diploma/internal/storage/auth"
	"golang.org/x/crypto/bcrypt"
)

const bcryptCost = 12

type Service struct {
	storage authStorage.Storage
}

func NewService(storage authStorage.Storage) (*Service, error) {
	if storage == nil {
		return nil, errors.New("storage required")
	}

	return &Service{storage}, nil
}

func (s *Service) Register(cred auth.Credentials, ctx context.Context) error {
	hash, errHash := bcrypt.GenerateFromPassword(cred.Password, bcryptCost)
	if errHash != nil {
		return fmt.Errorf("cannot hash password: %w", errHash)
	}

	u := auth.UserToCreate{
		Login:        cred.Login,
		PasswordHash: hash,
	}

	if _, errCreate := s.storage.CreateUser(u, ctx); errCreate != nil {
		return fmt.Errorf("cannot create user [%s]: %w", cred.Login, errCreate)
	}

	return nil
}

func (s *Service) Login(cred auth.Credentials, ctx context.Context) (auth.SignedUserID, error) {
	nilLogin := auth.SignedUserID{}

	u, errGet := s.storage.GetUserByLogin(cred.Login, ctx)
	if errGet != nil {
		return nilLogin, fmt.Errorf("cannot get user by sess [%s]: %w", cred.Login, errGet)
	}
	if u == nil {
		return nilLogin, fmt.Errorf("user not found by sess [%s]: %w", cred.Login, auth.ErrWrongCredentials)
	}

	if errValidate := bcrypt.CompareHashAndPassword(u.PasswordHash, cred.Password); errValidate != nil {
		return nilLogin, fmt.Errorf("password doesn't match for user [%s]: %w", cred.Login, auth.ErrWrongCredentials)
	}

	sigKey, errGenKey := generateKey()
	if errGenKey != nil {
		return nilLogin, fmt.Errorf("cannot generate signature key: %w", errGenKey)
	}

	sessToStart := auth.UserSessionToStart{
		UserID:       u.ID,
		SignatureKey: sigKey,
	}

	sess, errSet := s.storage.SetUserSession(sessToStart, ctx)
	if errSet != nil {
		return nilLogin, fmt.Errorf("cannot create session for user [%d]: %w", u.ID, errSet)
	}

	signedUserID := signUserID(sess)

	return signedUserID, nil
}

func (s *Service) Validate(sgn auth.SignedUserID, ctx context.Context) error {
	sess, errGet := s.storage.GetUserSession(sgn.ID, ctx)
	if errGet != nil {
		return fmt.Errorf("cannot get user session: %w", errGet)
	}

	canonicalS := signUserID(sess)
	if !hmac.Equal(canonicalS.Signature, sgn.Signature) {
		return fmt.Errorf("signature %v doesn't match", sgn)
	}

	return nil
}

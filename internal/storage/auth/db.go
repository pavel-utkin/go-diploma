package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	srv "go-diploma/internal/service/auth"
	"log"
)

type AuthStorage struct {
	db *sql.DB
}

func NewAuthStorage(db *sql.DB) (*AuthStorage, error) {
	if db == nil {
		return nil, errors.New("db should not be nil")
	}
	return &AuthStorage{db}, nil
}

func (s *AuthStorage) CreateUser(u srv.UserToCreate, ctx context.Context) (srv.User, error) {
	row := s.db.QueryRowContext(ctx, `
		insert into USERS (USERS_LOGIN, USERS_PASSWORD_HASH) 
		values($1, $2) 
		returning USERS_ID, USERS_LOGIN, USERS_PASSWORD_HASH
		`, u.Login, u.PasswordHash)
	user := srv.User{}

	err := row.Scan(&user.ID, &user.Login, &user.PasswordHash)
	var dbErr *pgconn.PgError
	if errors.As(err, &dbErr) && dbErr.Code == pgerrcode.UniqueViolation {
		log.Printf("Duplicate login: %s", u.Login)
		err = srv.ErrLoginAlreadyTaken
	}
	if err != nil {
		return user, fmt.Errorf("cannot insert user: %w", err)
	}

	return user, nil
}

func (s *AuthStorage) GetUserByLogin(login string, ctx context.Context) (*srv.User, error) {
	row := s.db.QueryRowContext(ctx, `
		select USERS_ID, USERS_LOGIN, USERS_PASSWORD_HASH
		from USERS
		where USERS_LOGIN = $1
		`, login)
	user := srv.User{}

	err := row.Scan(&user.ID, &user.Login, &user.PasswordHash)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, srv.ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("cannot read user from DB: %w", err)
	}

	return &user, nil
}

func (s *AuthStorage) SetUserSession(us srv.UserSessionToStart, ctx context.Context) (srv.UserSession, error) {
	row := s.db.QueryRowContext(ctx, `
		insert into USER_SESSIONS (USERS_ID, USER_SESSIONS_SIG_KEY) 
		values($1, $2) 
		on conflict (USERS_ID) do update set USER_SESSIONS_SIG_KEY = $2
		returning USERS_ID, USER_SESSIONS_SIG_KEY, USER_SESSIONS_STARTED_AT
		`, us.UserID, us.SignatureKey)
	sess := srv.UserSession{}

	if err := row.Scan(&sess.UserID, &sess.SignatureKey, &sess.StartedAt); err != nil {
		return sess, fmt.Errorf("cannot insert user session: %w", err)
	}

	return sess, nil
}

func (s *AuthStorage) GetUserSession(uID int64, ctx context.Context) (srv.UserSession, error) {
	row := s.db.QueryRowContext(ctx, `
		select USERS_ID, USER_SESSIONS_SIG_KEY, USER_SESSIONS_STARTED_AT
		from USER_SESSIONS 
		where USERS_ID = $1  
		`, uID)
	sess := srv.UserSession{}

	if err := row.Scan(&sess.UserID, &sess.SignatureKey, &sess.StartedAt); err != nil {
		return sess, fmt.Errorf("cannot get user session: %w", err)
	}

	return sess, nil
}

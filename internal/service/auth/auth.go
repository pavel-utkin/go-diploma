package auth

import "context"

type Service interface {
	Register(u Credentials, ctx context.Context) error
	Login(cred Credentials, ctx context.Context) (SignedUserID, error)
}

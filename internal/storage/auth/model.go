package auth

type UserToCreate struct {
	Login        string
	PasswordHash []byte
}

type User struct {
	ID           int64
	Login        string
	PasswordHash []byte
}

type UserSession struct {
	UserID       int64
	SignatureKey []byte
}

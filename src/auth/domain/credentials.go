package domain

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type (
	OwnerID string

	CredentialsID string

	Credentials struct {
		ID           CredentialsID
		OwnerID      OwnerID
		Email        string
		PasswordHash []byte
		CreatedAt    time.Time
		UpdatedAt    time.Time
	}
)

func (c *Credentials) VerifyPassword(password string) bool {
	return bcrypt.CompareHashAndPassword(c.PasswordHash, []byte(password)) == nil
}

func HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

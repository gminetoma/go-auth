package domain

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"time"

	"github.com/gminetoma/go-shared/src/id"
)

type (
	NewRefreshTokenParams struct {
		OwnerID   OwnerID
		ExpiresAt time.Time
		Now       time.Time
	}
)

func NewRefreshToken(params NewRefreshTokenParams) (*RefreshToken, string, error) {
	rawToken, err := generateToken(32)
	if err != nil {
		return nil, "", err
	}

	hash, err := HashRefreshToken(rawToken)
	if err != nil {
		return nil, "", err
	}

	return &RefreshToken{
		ID:        RefreshTokenID(id.Make()),
		OwnerID:   params.OwnerID,
		TokenHash: hash,
		ExpiresAt: params.ExpiresAt,
		CreatedAt: params.Now,
		RevokedAt: nil,
	}, rawToken, nil
}

func HashRefreshToken(token string) ([]byte, error) {
	decodedToken, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return nil, err
	}

	hash := sha256.Sum256(decodedToken)
	return hash[:], nil
}

func generateToken(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}

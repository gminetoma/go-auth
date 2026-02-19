package domain

import (
	"crypto/subtle"
	"time"
)

type (
	RefreshTokenID string

	RefreshToken struct {
		ID        RefreshTokenID
		OwnerID   OwnerID
		TokenHash []byte
		ExpiresAt time.Time
		CreatedAt time.Time
		RevokedAt *time.Time
	}
)

func (t *RefreshToken) Revoke(now time.Time) {
	t.RevokedAt = &now
}

func (t *RefreshToken) IsUsable(now time.Time, grace time.Duration) bool {
	if t.isExpired(now) {
		return false
	}

	if t.isRevoked() {
		return t.isWithinGracePeriod(now, grace)
	}

	return true
}

func (t *RefreshToken) VerifyToken(token string) bool {
	hash, err := HashRefreshToken(token)
	if err != nil {
		return false
	}

	return subtle.ConstantTimeCompare(t.TokenHash, hash) == 1
}

func (t *RefreshToken) isExpired(now time.Time) bool {
	return now.After(t.ExpiresAt)
}

func (t *RefreshToken) isWithinGracePeriod(now time.Time, grace time.Duration) bool {
	if t.RevokedAt == nil {
		return false
	}

	return t.RevokedAt.Add(grace).After(now)
}

func (t *RefreshToken) isRevoked() bool {
	return t.RevokedAt != nil
}

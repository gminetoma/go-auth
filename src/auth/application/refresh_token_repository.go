package application

import (
	"context"

	"github.com/gminetoma/go-auth/src/auth/domain"
)

type (
	RefreshTokenRepository interface {
		Create(ctx context.Context, refreshToken domain.RefreshToken) error
		FindByToken(ctx context.Context, token string) (*domain.RefreshToken, error)
		Update(ctx context.Context, refreshToken domain.RefreshToken) error
	}
)

package application

import (
	"context"

	"github.com/gminetoma/go-auth/src/auth/domain"
)

type (
	CredentialsRepository interface {
		Create(ctx context.Context, credentials domain.Credentials) error
		FindByEmail(ctx context.Context, email string) (*domain.Credentials, error)
	}
)

package application

import (
	"context"

	"github.com/gminetoma/go-auth/src/user/domain"
)

type (
	UserRepository interface {
		Create(ctx context.Context, user domain.User) error
		FindByID(ctx context.Context, id domain.UserID) (*domain.User, error)
	}
)

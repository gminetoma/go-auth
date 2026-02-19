package infrastructure

import (
	"context"
	"database/sql"
	"errors"

	"github.com/gminetoma/go-auth/src/auth/domain"
	"github.com/gminetoma/go-auth/src/auth/infrastructure/sqlc"
	"github.com/gminetoma/go-shared/src/errs"
)

type (
	PGCredentialsRepository struct {
		Queries *sqlc.Queries
	}
)

func (r *PGCredentialsRepository) Create(ctx context.Context, credentials domain.Credentials) error {
	return r.Queries.CreateCredentials(ctx, sqlc.CreateCredentialsParams{
		ID:           string(credentials.ID),
		OwnerID:      string(credentials.OwnerID),
		Email:        credentials.Email,
		PasswordHash: credentials.PasswordHash,
		CreatedAt:    credentials.CreatedAt,
		UpdatedAt:    credentials.UpdatedAt,
	})
}

func (r *PGCredentialsRepository) FindByEmail(ctx context.Context, email string) (*domain.Credentials, error) {
	credential, err := r.Queries.FindCredentialsByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNotFound
		}

		return nil, err
	}

	return &domain.Credentials{
		ID:           domain.CredentialsID(credential.ID),
		OwnerID:      domain.OwnerID(credential.OwnerID),
		Email:        credential.Email,
		PasswordHash: credential.PasswordHash,
		CreatedAt:    credential.CreatedAt,
		UpdatedAt:    credential.UpdatedAt,
	}, nil
}

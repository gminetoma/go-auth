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
	PGRefreshTokenRepository struct {
		Queries *sqlc.Queries
	}
)

func (r *PGRefreshTokenRepository) Create(ctx context.Context, refreshToken domain.RefreshToken) error {
	return r.Queries.CreateRefreshToken(ctx, sqlc.CreateRefreshTokenParams{
		ID:        string(refreshToken.ID),
		OwnerID:   string(refreshToken.OwnerID),
		TokenHash: refreshToken.TokenHash,
		ExpiresAt: refreshToken.ExpiresAt,
		CreatedAt: refreshToken.CreatedAt,
	})
}

func (r *PGRefreshTokenRepository) FindByToken(ctx context.Context, token string) (*domain.RefreshToken, error) {
	hash, err := domain.HashRefreshToken(token)
	if err != nil {
		return nil, err
	}

	rt, err := r.Queries.FindRefreshTokenByToken(ctx, hash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNotFound
		}

		return nil, err
	}

	return &domain.RefreshToken{
		ID:        domain.RefreshTokenID(rt.ID),
		OwnerID:   domain.OwnerID(rt.OwnerID),
		TokenHash: rt.TokenHash,
		ExpiresAt: rt.ExpiresAt,
		CreatedAt: rt.CreatedAt,
	}, nil
}

func (r *PGRefreshTokenRepository) Update(ctx context.Context, refreshToken domain.RefreshToken) error {
	var revokedAt sql.NullTime
	if refreshToken.RevokedAt != nil {
		revokedAt = sql.NullTime{Time: *refreshToken.RevokedAt, Valid: true}
	}

	rows, err := r.Queries.UpdateRefreshToken(ctx, sqlc.UpdateRefreshTokenParams{
		ID:        string(refreshToken.ID),
		RevokedAt: revokedAt,
	})
	if err != nil {
		return err
	}

	if rows == 0 {
		return errs.ErrNotFound
	}

	return nil
}

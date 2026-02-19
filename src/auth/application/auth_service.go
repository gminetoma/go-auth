package application

import (
	"context"
	"errors"
	"time"

	"github.com/gminetoma/go-auth/src/auth/domain"
	"github.com/gminetoma/go-auth/src/auth/infrastructure"
	"github.com/gminetoma/go-auth/src/auth/infrastructure/sqlc"
	"github.com/gminetoma/go-shared/src/clock"
	"github.com/gminetoma/go-shared/src/errs"
	"github.com/gminetoma/go-shared/src/identity"
)

type (
	TokenPair struct {
		AccessToken  string
		RefreshToken string
	}

	AuthService interface {
		Login(ctx context.Context, email, password string) (*TokenPair, error)
		Refresh(ctx context.Context, refreshToken string) (*TokenPair, error)
		Logout(ctx context.Context, refreshToken string) error
		Me(ctx context.Context) (domain.OwnerID, error)
		VerifyAccessToken(token string) (domain.OwnerID, error)
		Register(ctx context.Context, email, password, ownerID string) error
	}

	authService struct {
		RefreshTokenRepository  RefreshTokenRepository
		CredentialsRepository   CredentialsRepository
		TokenService            TokenService
		RefreshTokenExpiry      time.Duration
		RefreshTokenGracePeriod time.Duration
		Clock                   clock.Clock
	}
)

type (
	NewAuthParams struct {
		DB                      sqlc.DBTX
		AccessTokenSecret       string
		RefreshTokenExpiry      time.Duration
		AccessTokenExpiry       time.Duration
		RefreshTokenGracePeriod time.Duration
		Clock                   clock.Clock
	}
)

func NewAuthService(params NewAuthParams) AuthService {
	queries := sqlc.New(params.DB)

	return &authService{
		RefreshTokenRepository: &infrastructure.PGRefreshTokenRepository{
			Queries: queries,
		},
		CredentialsRepository: &infrastructure.PGCredentialsRepository{
			Queries: queries,
		},
		TokenService: &infrastructure.JWTTokenService{
			Secret: []byte(params.AccessTokenSecret),
			Expiry: params.AccessTokenExpiry,
		},
		RefreshTokenExpiry:      params.RefreshTokenExpiry,
		RefreshTokenGracePeriod: params.RefreshTokenGracePeriod,
		Clock:                   params.Clock,
	}
}

func (s *authService) Register(ctx context.Context, email, password, ownerID string) error {
	creds, errs := domain.NewCredentials(domain.NewCredentialsParams{
		OwnerID:  domain.OwnerID(ownerID),
		Email:    email,
		Password: password,
		Now:      s.Clock.Now(),
	})
	if errs != nil {
		return errors.Join(errs...)
	}

	return s.CredentialsRepository.Create(ctx, *creds)
}

func (s *authService) Login(ctx context.Context, email, password string) (*TokenPair, error) {
	creds, err := s.CredentialsRepository.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return nil, ErrInvalidCredentials
		}

		return nil, err
	}

	if !creds.VerifyPassword(password) {
		return nil, ErrInvalidCredentials
	}

	return s.generateTokenPair(ctx, creds.OwnerID)
}

func (s *authService) Refresh(ctx context.Context, refreshToken string) (*TokenPair, error) {

	existing, err := s.revokeRefreshTokenByToken(ctx, refreshToken, s.RefreshTokenGracePeriod)
	if err != nil {
		return nil, err
	}

	if err := s.RefreshTokenRepository.Update(ctx, *existing); err != nil {
		return nil, err
	}

	return s.generateTokenPair(ctx, existing.OwnerID)
}

func (s *authService) Logout(ctx context.Context, refreshToken string) error {
	if _, err := identity.Require(ctx); err != nil {
		return err
	}

	existing, err := s.revokeRefreshTokenByToken(ctx, refreshToken, s.RefreshTokenGracePeriod)
	if err != nil {
		return err
	}

	return s.RefreshTokenRepository.Update(ctx, *existing)
}

func (s *authService) Me(ctx context.Context) (domain.OwnerID, error) {
	ownerID, err := identity.Require(ctx)
	if err != nil {
		return "", err
	}

	return domain.OwnerID(ownerID), nil
}

func (s *authService) VerifyAccessToken(token string) (domain.OwnerID, error) {
	return s.TokenService.Verify(token)
}

func (s *authService) revokeRefreshTokenByToken(ctx context.Context, token string, refreshTokenGracePeriod time.Duration) (*domain.RefreshToken, error) {
	existing, err := s.RefreshTokenRepository.FindByToken(ctx, token)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return nil, ErrInvalidRefreshToken
		}

		return nil, err
	}

	if !existing.IsUsable(s.Clock.Now(), refreshTokenGracePeriod) {
		return nil, ErrExpiredRefreshToken
	}

	existing.Revoke(s.Clock.Now())

	return existing, nil
}

func (s *authService) generateTokenPair(ctx context.Context, ownerID domain.OwnerID) (*TokenPair, error) {
	accessToken, err := s.TokenService.Generate(ownerID)
	if err != nil {
		return nil, err
	}

	now := s.Clock.Now()

	refreshToken, rawRefreshToken, err := domain.NewRefreshToken(domain.NewRefreshTokenParams{
		OwnerID:   ownerID,
		ExpiresAt: now.Add(s.RefreshTokenExpiry),
		Now:       now,
	})
	if err != nil {
		return nil, err
	}

	if err := s.RefreshTokenRepository.Create(ctx, *refreshToken); err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: rawRefreshToken,
	}, nil
}

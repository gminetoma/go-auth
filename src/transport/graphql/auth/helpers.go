package auth

import (
	"context"

	"github.com/gminetoma/go-auth/src/auth/application"
	"github.com/gminetoma/go-auth/src/transport/cookies"
	"github.com/gminetoma/go-auth/src/transport/middleware"
)

func RefreshTokenFromContext(ctx context.Context, token *string) (string, error) {
	if token != nil && *token != "" {
		return *token, nil
	}

	req, ok := middleware.Request(ctx)
	if !ok {
		return "", application.ErrRequiredRefreshToken
	}

	if t, found := cookies.RefreshToken(req); found {
		return t, nil
	}

	return "", application.ErrRequiredRefreshToken
}

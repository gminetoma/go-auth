package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gminetoma/go-auth/src/transport/cookies"
	"github.com/gminetoma/go-shared/src/identity"
)

type (
	TokenPair struct {
		AccessToken  string
		RefreshToken string
	}

	AuthParams struct {
		VerifyAccessToken  func(token string) (string, error)
		RefreshTokens      func(ctx context.Context, token string) (*TokenPair, error)
		AccessTokenExpiry  time.Duration
		RefreshTokenExpiry time.Duration
	}
)

func Auth(params AuthParams) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			token := accessTokenFromRequest(r)

			ownerID, err := params.VerifyAccessToken(token)
			if err == nil {
				ctx := identity.SetOwnerID(r.Context(), ownerID)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			// Access token invalid/missing - try auto-refresh
			refreshToken, found := cookies.RefreshToken(r)
			if !found {
				// No refresh token - continue without auth
				next.ServeHTTP(w, r)
				return
			}

			// Try to refresh tokens
			pair, err := params.RefreshTokens(r.Context(), refreshToken)
			if err != nil {
				// Refresh failed - continue without auth
				// (Could be expired, already used, or deleted)
				next.ServeHTTP(w, r)
				return
			}

			// Refresh successful - set new cookies
			cookies.SetAccessToken(w, pair.AccessToken, params.AccessTokenExpiry)
			cookies.SetRefreshToken(w, pair.RefreshToken, params.RefreshTokenExpiry)

			// Verify new access token to get owner ID
			ownerID, err = params.VerifyAccessToken(pair.AccessToken)
			if err != nil {
				// This shouldn't happen - Something is wrong with the token issuer
				slog.Error("Failed to verify refreshed access token", "error", err)

				cookies.ClearAuth(w)

				next.ServeHTTP(w, r)
				return
			}

			ctx := identity.SetOwnerID(r.Context(), ownerID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func accessTokenFromHeader(r *http.Request) string {
	header := r.Header.Get("Authorization")
	if header == "" {
		return ""
	}

	token := strings.TrimPrefix(header, "Bearer ")
	if token == header {
		return ""
	}

	return token
}

func accessTokenFromCookie(r *http.Request) string {
	cookie, err := r.Cookie(cookies.AccessTokenCookie)
	if err != nil {
		return ""
	}

	return cookie.Value
}

func accessTokenFromRequest(r *http.Request) string {
	token := accessTokenFromHeader(r)
	if token == "" {
		return accessTokenFromCookie(r)
	}

	return token
}

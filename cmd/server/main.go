package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	authApplication "github.com/gminetoma/go-auth/src/auth/application"
	"github.com/gminetoma/go-auth/src/transport/graphql"
	"github.com/gminetoma/go-auth/src/transport/middleware"
	"github.com/gminetoma/go-shared/src/clock"
	"github.com/gminetoma/go-shared/src/config"
	"github.com/gminetoma/go-shared/src/database"
)

func main() {
	cfg := config.LoadEnv()

	// Database
	db := database.PGXConnect(cfg.DatabaseURL)
	defer db.Close()

	mux := http.NewServeMux()

	clock := clock.New()

	authService := authApplication.NewAuthService(authApplication.NewAuthParams{
		DB:                      db,
		AccessTokenSecret:       cfg.JWTSecret,
		RefreshTokenExpiry:      cfg.RefreshTokenExpiry,
		AccessTokenExpiry:       cfg.AccessTokenExpiry,
		RefreshTokenGracePeriod: cfg.RefreshTokenGracePeriod,
		Clock:                   clock,
	})

	graphql.SetupGraphQL(mux, graphql.SetupGraphQLParams{
		AuthService:         authService,
		EnablePlayground:    cfg.Env != config.EnvProduction,
		EnableIntrospection: cfg.Env != config.EnvProduction,
	})

	// Middleware
	handler := middleware.Setup(
		mux,
		middleware.HTTPContext,
		middleware.Auth(middleware.AuthParams{
			VerifyAccessToken: func(token string) (string, error) {
				ownerID, err := authService.VerifyAccessToken(token)
				if err != nil {
					return "", err
				}

				return string(ownerID), nil
			},
			RefreshTokens: func(ctx context.Context, token string) (*middleware.TokenPair, error) {
				tokenPair, err := authService.Refresh(ctx, token)
				if err != nil {
					return nil, err
				}

				return &middleware.TokenPair{
					AccessToken:  tokenPair.AccessToken,
					RefreshToken: tokenPair.RefreshToken,
				}, err
			},
			AccessTokenExpiry:  cfg.AccessTokenExpiry,
			RefreshTokenExpiry: cfg.RefreshTokenExpiry,
		}),
	)

	// Server

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      http.TimeoutHandler(handler, 30*time.Second, "request timeout"),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 35 * time.Second,
	}

	serverAddr := fmt.Sprintf("%s://%s:%s", cfg.Protocol, cfg.Host, cfg.Port)
	slog.Info("server started", "address", serverAddr)

	if err := server.ListenAndServe(); err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}

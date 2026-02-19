package main

import (
	"context"
	"errors"
	"log/slog"
	"os"

	authDomain "github.com/gminetoma/go-auth/src/auth/domain"
	authInfrastructure "github.com/gminetoma/go-auth/src/auth/infrastructure"
	authSQLC "github.com/gminetoma/go-auth/src/auth/infrastructure/sqlc"
	userDomain "github.com/gminetoma/go-auth/src/user/domain"
	userInfrastructure "github.com/gminetoma/go-auth/src/user/infrastructure"
	userSQLC "github.com/gminetoma/go-auth/src/user/infrastructure/sqlc"
	"github.com/gminetoma/go-shared/src/clock"
	"github.com/gminetoma/go-shared/src/config"
	"github.com/gminetoma/go-shared/src/database"
	"github.com/gminetoma/go-shared/src/errs"
)

const (
	adminEmail    = "admin@example.com"
	adminPassword = "Password@123"
)

func main() {
	cfg := config.LoadEnv()

	db := database.PGXConnect(cfg.DatabaseURL)
	defer db.Close()

	ctx := context.Background()

	userRepo := &userInfrastructure.PGUserRepository{
		Queries: userSQLC.New(db),
	}

	credsRepo := &authInfrastructure.PGCredentialsRepository{
		Queries: authSQLC.New(db),
	}

	existing, err := credsRepo.FindByEmail(ctx, adminEmail)
	if err != nil {
		if !errors.Is(errs.ErrNotFound, err) {
			slog.Error("failed to check existing credentials", "error", err)
			os.Exit(1)
		}
	}

	if existing != nil {
		slog.Info("seed already exists, skipping", "email", adminEmail, "password", adminPassword)
		return
	}

	c := clock.New()

	user := userDomain.NewUser(userDomain.NewUserParams{
		Now: c.Now(),
	})

	if err := userRepo.Create(ctx, *user); err != nil {
		slog.Error("failed to create user", "error", err)
		os.Exit(1)
	}

	creds, errs := authDomain.NewCredentials(authDomain.NewCredentialsParams{
		OwnerID:  authDomain.OwnerID(user.ID),
		Email:    adminEmail,
		Password: adminPassword,
		Now:      c.Now(),
	})
	if errs != nil {
		slog.Error("failed to create credentials", "error", err)
		os.Exit(1)
	}

	if err := credsRepo.Create(ctx, *creds); err != nil {
		slog.Error("failed to persist credentials", "error", err)
		os.Exit(1)
	}

	slog.Info("seed completed", "email", adminEmail, "password", adminPassword)
}

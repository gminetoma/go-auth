package application

import "errors"

var (
	ErrInvalidCredentials   = errors.New("auth.invalid-credentials")
	ErrInvalidRefreshToken  = errors.New("auth.invalid-refresh-token")
	ErrExpiredRefreshToken  = errors.New("auth.expired-refresh-token")
	ErrRequiredRefreshToken = errors.New("auth.required-refresh-token")
)

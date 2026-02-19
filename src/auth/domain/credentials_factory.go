package domain

import (
	"time"

	"github.com/gminetoma/go-shared/src/id"
	"github.com/gminetoma/go-shared/src/validation"
)

type (
	NewCredentialsParams struct {
		OwnerID  OwnerID
		Email    string
		Password string
		Now      time.Time
	}
)

const passwordMinLength = 8
const passwordMaxLength = 72

func NewCredentials(params NewCredentialsParams) (*Credentials, []error) {
	validator := validation.New()

	validator.RequiredString(string(params.OwnerID), ErrOwnerIDRequired)

	validator.RequiredString(params.Email, ErrEmailRequired)
	validator.ValidEmail(params.Email, ErrEmailInvalid)

	validator.RequiredString(params.Password, ErrPasswordRequired)
	validator.MinLength(params.Password, passwordMinLength, ErrPasswordTooShort)
	validator.MaxLength(params.Password, passwordMaxLength, ErrPasswordTooLong)

	if validator.HasErrors() {
		return nil, validator.Errors()
	}

	passwordHash, err := HashPassword(params.Password)
	if err != nil {
		return nil, []error{err}
	}

	return &Credentials{
		ID:           CredentialsID(id.Make()),
		OwnerID:      params.OwnerID,
		Email:        params.Email,
		PasswordHash: passwordHash,
		CreatedAt:    params.Now,
		UpdatedAt:    params.Now,
	}, nil
}

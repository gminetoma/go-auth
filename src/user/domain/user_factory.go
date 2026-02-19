package domain

import (
	"time"

	"github.com/gminetoma/go-shared/src/id"
)

type (
	NewUserParams struct {
		Now time.Time
	}
)

func NewUser(params NewUserParams) *User {
	return &User{
		ID:        UserID(id.Make()),
		CreatedAt: params.Now,
	}
}

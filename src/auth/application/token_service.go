package application

import "github.com/gminetoma/go-auth/src/auth/domain"

type TokenService interface {
	Generate(ownerID domain.OwnerID) (string, error)
	Verify(token string) (domain.OwnerID, error)
}

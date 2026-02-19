package infrastructure

import (
	"fmt"
	"time"

	"github.com/gminetoma/go-auth/src/auth/domain"

	"github.com/golang-jwt/jwt/v5"
)

type (
	JWTTokenService struct {
		Secret []byte
		Expiry time.Duration
	}
)

func (s *JWTTokenService) Generate(ownerID domain.OwnerID) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   string(ownerID),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.Expiry)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(s.Secret)
}

func (s *JWTTokenService) Verify(tokenString string) (domain.OwnerID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(t *jwt.Token) (any, error) {
		// Prevent algorithm confusion attacks where an attacker changes the token's
		// "alg" header to bypass signature verification.
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return s.Secret, nil
	})
	if err != nil {
		return "", err
	}

	subject, err := token.Claims.GetSubject()
	if err != nil {
		return "", err
	}

	return domain.OwnerID(subject), nil
}

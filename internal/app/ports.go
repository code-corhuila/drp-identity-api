package app

import (
	"context"

	"github.com/code-corhuila/drp-identity-api/internal/domain"
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email domain.Email) (domain.User, error)
	FindByID(ctx context.Context, id string) (domain.User, error)
}

type PasswordChecker interface {
	Match(hash, password string) bool
}

type TokenService interface {
	Issue(user domain.User) (token string, expiresIn int, err error)
	PublicJWKS() map[string]any
	Parse(token string) (userID string, err error)
}

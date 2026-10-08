package app

import (
	"context"
	"strings"

	"github.com/code-corhuila/drp-identity-api/internal/domain"
)

type LoginCommand struct {
	Email    string
	Password string
}

type LoginResult struct {
	AccessToken string
	TokenType   string
	ExpiresIn   int
	User        domain.User
}

type Auth struct {
	Users     UserRepository
	Passwords PasswordChecker
	Tokens    TokenService
}

func (a Auth) Login(ctx context.Context, cmd LoginCommand) (LoginResult, error) {
	if strings.TrimSpace(cmd.Email) == "" || len(cmd.Password) < 8 {
		return LoginResult{}, ErrValidation
	}
	email, err := domain.NewEmail(cmd.Email)
	if err != nil {
		return LoginResult{}, ErrValidation
	}
	user, err := a.Users.FindByEmail(ctx, email)
	if err != nil || user.DeletedAt != nil {
		return LoginResult{}, ErrInvalidCredentials
	}
	if !a.Passwords.Match(user.PasswordHash, cmd.Password) {
		return LoginResult{}, ErrInvalidCredentials
	}
	tok, exp, err := a.Tokens.Issue(user)
	if err != nil {
		return LoginResult{}, err
	}
	return LoginResult{AccessToken: tok, TokenType: "Bearer", ExpiresIn: exp, User: user}, nil
}

func (a Auth) CurrentUser(ctx context.Context, bearer string) (domain.User, error) {
	raw := strings.TrimSpace(bearer)
	raw = strings.TrimPrefix(raw, "Bearer ")
	raw = strings.TrimPrefix(raw, "bearer ")
	if raw == "" {
		return domain.User{}, ErrUnauthorized
	}
	id, err := a.Tokens.Parse(raw)
	if err != nil {
		return domain.User{}, ErrUnauthorized
	}
	user, err := a.Users.FindByID(ctx, id)
	if err != nil || user.DeletedAt != nil {
		return domain.User{}, ErrUnauthorized
	}
	return user, nil
}

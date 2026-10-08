package app

import (
	"context"
	"testing"
	"time"

	"github.com/code-corhuila/drp-identity-api/internal/domain"
)

type stubUsers struct{ u domain.User }

func (s stubUsers) FindByEmail(_ context.Context, email domain.Email) (domain.User, error) {
	if email != s.u.Email {
		return domain.User{}, ErrUserNotFound
	}
	return s.u, nil
}

func (s stubUsers) FindByID(_ context.Context, id string) (domain.User, error) {
	if id != s.u.ID {
		return domain.User{}, ErrUserNotFound
	}
	return s.u, nil
}

type stubPwd struct{ ok bool }

func (s stubPwd) Match(_, _ string) bool { return s.ok }

type stubTok struct{}

func (stubTok) Issue(u domain.User) (string, int, error) { return "tok-"+u.ID, 3600, nil }
func (stubTok) PublicJWKS() map[string]any               { return map[string]any{"keys": []any{}} }
func (stubTok) Parse(token string) (string, error) {
	if token == "bad" {
		return "", ErrUnauthorized
	}
	return stringsTrimPrefix(token, "tok-"), nil
}

func stringsTrimPrefix(s, p string) string {
	if len(s) >= len(p) && s[:len(p)] == p {
		return s[len(p):]
	}
	return s
}

func testUser(t *testing.T) domain.User {
	t.Helper()
	u, err := domain.NewUser("44444444-4444-4444-4444-444444444444", "member@spacehub.local", "hash", "Ana Reserva", []domain.Role{domain.RoleUser}, time.Unix(0, 0).UTC())
	if err != nil {
		t.Fatal(err)
	}
	return u
}

func TestLoginOK(t *testing.T) {
	u := testUser(t)
	a := Auth{Users: stubUsers{u}, Passwords: stubPwd{true}, Tokens: stubTok{}}
	got, err := a.Login(context.Background(), LoginCommand{Email: "Member@spacehub.local", Password: "Spacehub1!"})
	if err != nil || got.AccessToken == "" || got.TokenType != "Bearer" {
		t.Fatalf("%v %+v", err, got)
	}
}

func TestLoginRejects(t *testing.T) {
	u := testUser(t)
	a := Auth{Users: stubUsers{u}, Passwords: stubPwd{false}, Tokens: stubTok{}}
	if _, err := a.Login(context.Background(), LoginCommand{Email: "x", Password: "short"}); err != ErrValidation {
		t.Fatalf("short password: %v", err)
	}
	if _, err := a.Login(context.Background(), LoginCommand{Email: "member@spacehub.local", Password: ""}); err != ErrInvalidCredentials {
		t.Fatalf("empty password: %v", err)
	}
	if _, err := a.Login(context.Background(), LoginCommand{Email: "member@spacehub.local", Password: "Spacehub1!"}); err != ErrInvalidCredentials {
		t.Fatalf("bad password: %v", err)
	}
	unknown := Auth{Users: stubUsers{u}, Passwords: stubPwd{true}, Tokens: stubTok{}}
	if _, err := unknown.Login(context.Background(), LoginCommand{Email: "nobody@spacehub.local", Password: "Spacehub1!"}); err != ErrInvalidCredentials {
		t.Fatalf("unknown user: %v", err)
	}
}

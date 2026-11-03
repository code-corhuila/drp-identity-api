package security

import (
	"testing"
	"time"

	"github.com/code-corhuila/drp-identity-api/internal/domain"
	"github.com/golang-jwt/jwt/v5"
)

func TestIssueHasE01Claims(t *testing.T) {
	s, err := LoadOrGenerate("", "test-kid")
	if err != nil {
		t.Fatal(err)
	}
	u, err := domain.NewUser("44444444-4444-4444-4444-444444444444", "member@spacehub.local", "hash", "Ana", []domain.Role{domain.RoleUser}, time.Unix(0, 0).UTC())
	if err != nil {
		t.Fatal(err)
	}
	raw, exp, err := s.Issue(u)
	if err != nil || exp != 3600 {
		t.Fatalf("issue %v exp %d", err, exp)
	}
	tok, err := jwt.Parse(raw, func(t *jwt.Token) (any, error) {
		return &s.key.PublicKey, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Alg()}))
	if err != nil {
		t.Fatal(err)
	}
	claims := tok.Claims.(jwt.MapClaims)
	if claims["sub"] != u.ID || claims["role"] != "USER" || claims["jti"] == nil || claims["jti"] == "" {
		t.Fatalf("claims %+v", claims)
	}
}

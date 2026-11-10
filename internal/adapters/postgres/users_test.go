package postgres

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/code-corhuila/drp-identity-api/internal/domain"
)

func TestFindByEmailLiveCorte2(t *testing.T) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		t.Skip("DATABASE_URL unset; Flyway identity schema not required in unit CI")
	}
	db, err := Open(dsn)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	users := NewUsers(db)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := users.EnsureSeedPassword(ctx, "test-seed-only"); err != nil {
		t.Fatal(err)
	}
	email, err := domain.NewEmail("member@spacehub.local")
	if err != nil {
		t.Fatal(err)
	}
	got, err := users.FindByEmail(ctx, email)
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != "44444444-4444-4444-4444-444444444444" || !got.HasRole(domain.RoleUser) {
		t.Fatalf("got %+v", got)
	}
}

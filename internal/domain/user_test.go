package domain

import (
	"testing"
	"time"
)

func TestNewUser(t *testing.T) {
	now := time.Date(2026, 10, 6, 12, 0, 0, 0, time.UTC)
	u, err := NewUser("44444444-4444-4444-4444-444444444444", "Member@spacehub.local", "hash", "Ana", []Role{RoleUser}, now)
	if err != nil {
		t.Fatal(err)
	}
	if u.Email.String() != "member@spacehub.local" {
		t.Fatalf("email not normalized: %s", u.Email)
	}
	if !u.HasRole(RoleUser) || u.HasRole(RoleAdmin) {
		t.Fatalf("roles %+v", u.Roles)
	}
}

func TestNewUserRejects(t *testing.T) {
	now := time.Now().UTC()
	cases := []struct {
		email, hash, name string
		roles             []Role
	}{
		{"", "h", "Ana", []Role{RoleUser}},
		{"a@b.c", "", "Ana", []Role{RoleUser}},
		{"a@b.c", "h", "", []Role{RoleUser}},
		{"a@b.c", "h", "Ana", nil},
		{"a@b.c", "h", "Ana", []Role{"GUEST"}},
	}
	for _, c := range cases {
		if _, err := NewUser("id", c.email, c.hash, c.name, c.roles, now); err == nil {
			t.Fatalf("expected error for %+v", c)
		}
	}
}

func TestSoftDelete(t *testing.T) {
	now := time.Now().UTC()
	u, _ := NewUser("id", "a@b.c", "h", "Ana", []Role{RoleUser}, now)
	u, err := u.SoftDelete(now.Add(time.Second))
	if err != nil || u.DeletedAt == nil {
		t.Fatalf("err=%v deleted=%v", err, u.DeletedAt)
	}
	if _, err := u.SoftDelete(now.Add(2 * time.Second)); err == nil {
		t.Fatal("second delete should fail")
	}
}

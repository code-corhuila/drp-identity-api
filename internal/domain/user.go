package domain

import (
	"strings"
	"time"
)

type Role string

const (
	RoleUser  Role = "USER"
	RoleAdmin Role = "ADMIN"
)

func ParseRole(v string) (Role, error) {
	switch Role(v) {
	case RoleUser, RoleAdmin:
		return Role(v), nil
	default:
		return "", ErrInvalidInput
	}
}

type Email string

func NewEmail(raw string) (Email, error) {
	e := strings.ToLower(strings.TrimSpace(raw))
	if e == "" || len(e) > 320 || !strings.Contains(e, "@") {
		return "", ErrInvalidInput
	}
	return Email(e), nil
}

func (e Email) String() string { return string(e) }

// User is the identity aggregate root. Password hashes live here; other
// services store only userId (BR-007 / BR-008).
type User struct {
	ID           string
	Email        Email
	PasswordHash string
	DisplayName  string
	Roles        []Role
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}

func NewUser(id, email, passwordHash, displayName string, roles []Role, now time.Time) (User, error) {
	em, err := NewEmail(email)
	if err != nil {
		return User{}, err
	}
	name := strings.TrimSpace(displayName)
	if id == "" || passwordHash == "" || name == "" || len(name) > 200 {
		return User{}, ErrInvalidInput
	}
	if err := validateRoles(roles); err != nil {
		return User{}, err
	}
	return User{
		ID:           id,
		Email:        em,
		PasswordHash: passwordHash,
		DisplayName:  name,
		Roles:        append([]Role(nil), roles...),
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

func validateRoles(roles []Role) error {
	if len(roles) == 0 {
		return ErrInvalidInput
	}
	seen := map[Role]struct{}{}
	for _, r := range roles {
		if _, err := ParseRole(string(r)); err != nil {
			return err
		}
		if _, ok := seen[r]; ok {
			return ErrInvalidInput
		}
		seen[r] = struct{}{}
	}
	return nil
}

func (u User) HasRole(r Role) bool {
	for _, x := range u.Roles {
		if x == r {
			return true
		}
	}
	return false
}

func (u User) SoftDelete(now time.Time) (User, error) {
	if u.DeletedAt != nil {
		return User{}, ErrInvalidTransition
	}
	t := now
	u.DeletedAt = &t
	u.UpdatedAt = now
	return u, nil
}

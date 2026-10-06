package memory

import (
	"context"
	"sync"
	"time"

	"github.com/code-corhuila/drp-identity-api/internal/app"
	"github.com/code-corhuila/drp-identity-api/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

// Users is an in-memory port. Flyway persistence is a later feat.
// Corte 2 UI does not read this: drp-front stays on synthetic/failover.
type Users struct {
	mu   sync.RWMutex
	byID map[string]domain.User
	mail map[domain.Email]string
}

func Corte2(plainPassword string) (*Users, error) {
	if plainPassword == "" {
		plainPassword = "Spacehub1!"
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.MinCost)
	if err != nil {
		return nil, err
	}
	now := time.Unix(0, 0).UTC()
	member, err := domain.NewUser("44444444-4444-4444-4444-444444444444", "member@spacehub.local", string(hash), "Ana Reserva", []domain.Role{domain.RoleUser}, now)
	if err != nil {
		return nil, err
	}
	admin, err := domain.NewUser("66666666-6666-6666-6666-666666666666", "admin@spacehub.local", string(hash), "Admin SpaceHub", []domain.Role{domain.RoleAdmin}, now)
	if err != nil {
		return nil, err
	}
	r := &Users{byID: map[string]domain.User{}, mail: map[domain.Email]string{}}
	r.put(member)
	r.put(admin)
	return r, nil
}

func (r *Users) put(u domain.User) {
	r.byID[u.ID] = u
	r.mail[u.Email] = u.ID
}

func (r *Users) FindByEmail(_ context.Context, email domain.Email) (domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	id, ok := r.mail[email]
	if !ok {
		return domain.User{}, app.ErrUserNotFound
	}
	return r.byID[id], nil
}

func (r *Users) FindByID(_ context.Context, id string) (domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, ok := r.byID[id]
	if !ok {
		return domain.User{}, app.ErrUserNotFound
	}
	return u, nil
}

package postgres

import (
	"context"
	"database/sql"
	"strings"

	"github.com/code-corhuila/drp-identity-api/internal/app"
	"github.com/code-corhuila/drp-identity-api/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

// Users reads identity.users / identity.user_roles. No DDL here (Anexo J).
type Users struct {
	DB *sql.DB
}

func NewUsers(db *sql.DB) *Users {
	return &Users{DB: db}
}

func (u *Users) FindByEmail(ctx context.Context, email domain.Email) (domain.User, error) {
	return u.load(ctx, `u.email = $1`, string(email))
}

func (u *Users) FindByID(ctx context.Context, id string) (domain.User, error) {
	return u.load(ctx, `u.id = $1`, id)
}

func (u *Users) load(ctx context.Context, where string, arg any) (domain.User, error) {
	rows, err := u.DB.QueryContext(ctx, `
SELECT u.id, u.email, u.password_hash, u.display_name, u.created_at, u.updated_at, u.deleted_at, r.role
FROM users u
LEFT JOIN user_roles r ON r.user_id = u.id
WHERE `+where+` AND u.deleted_at IS NULL
ORDER BY r.role`, arg)
	if err != nil {
		return domain.User{}, err
	}
	defer rows.Close()
	var got domain.User
	var roles []domain.Role
	found := false
	for rows.Next() {
		var deleted sql.NullTime
		var role sql.NullString
		var email string
		if err := rows.Scan(&got.ID, &email, &got.PasswordHash, &got.DisplayName, &got.CreatedAt, &got.UpdatedAt, &deleted, &role); err != nil {
			return domain.User{}, err
		}
		em, err := domain.NewEmail(email)
		if err != nil {
			return domain.User{}, err
		}
		got.Email = em
		if deleted.Valid {
			t := deleted.Time.UTC()
			got.DeletedAt = &t
		}
		got.CreatedAt = got.CreatedAt.UTC()
		got.UpdatedAt = got.UpdatedAt.UTC()
		if role.Valid {
			parsed, perr := domain.ParseRole(role.String)
			if perr == nil {
				roles = append(roles, parsed)
			}
		}
		found = true
	}
	if err := rows.Err(); err != nil {
		return domain.User{}, err
	}
	if !found {
		return domain.User{}, app.ErrUserNotFound
	}
	got.Roles = roles
	return got, nil
}

// EnsureSeedPassword replaces Flyway placeholder hashes so login works (DML only).
func (u *Users) EnsureSeedPassword(ctx context.Context, plain string) error {
	if strings.TrimSpace(plain) == "" {
		plain = "Spacehub1!"
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = u.DB.ExecContext(ctx, `
UPDATE users SET password_hash = $1, updated_at = NOW()
WHERE password_hash LIKE 'PENDING%'`, string(hash))
	return err
}

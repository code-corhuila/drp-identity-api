package security

import "golang.org/x/crypto/bcrypt"

type Passwords struct{}

func (Passwords) Match(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

package app

import "errors"

var (
	ErrValidation          = errors.New("validation")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrUserNotFound        = errors.New("user not found")
)

package httpadapter

import (
	"net/http"

	"github.com/code-corhuila/drp-identity-api/internal/app"
)

type Deps struct {
	Service string
	Auth    app.Auth
}

func NewMux(d Deps) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("GET /health", HealthHandler(d.Service))
	if d.Auth.Tokens != nil {
		mux.Handle("GET /api/v1/auth/jwks", JWKSHandler(d.Auth.Tokens))
	}
	if d.Auth.Users != nil {
		mux.Handle("POST /api/v1/auth/login", LoginHandler(d.Auth))
		mux.Handle("GET /api/v1/users/me", MeHandler(d.Auth))
	}
	mux.Handle("/", NotFoundHandler())
	return WithCorrelation(stripSpoofedIdentity(mux))
}

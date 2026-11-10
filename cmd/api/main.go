package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	httpadapter "github.com/code-corhuila/drp-identity-api/internal/adapters/http"
	"github.com/code-corhuila/drp-identity-api/internal/adapters/memory"
	"github.com/code-corhuila/drp-identity-api/internal/adapters/postgres"
	"github.com/code-corhuila/drp-identity-api/internal/adapters/security"
	"github.com/code-corhuila/drp-identity-api/internal/app"
)

func main() {
	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":8081"
	}
	service := os.Getenv("SERVICE_NAME")
	if service == "" {
		service = "identity-service"
	}

	var users app.UserRepository
	store := "memory Corte 2 seed"
	if dsn := os.Getenv("DATABASE_URL"); dsn != "" {
		db, err := postgres.Open(dsn)
		if err != nil {
			log.Fatal(err)
		}
		db.SetMaxOpenConns(8)
		pg := postgres.NewUsers(db)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := pg.EnsureSeedPassword(ctx, os.Getenv("IDENTITY_SEED_PASSWORD")); err != nil {
			log.Fatal(err)
		}
		users = pg
		store = "postgres identity schema"
	} else {
		seed, err := memory.Corte2(os.Getenv("IDENTITY_SEED_PASSWORD"))
		if err != nil {
			log.Fatal(err)
		}
		users = seed
	}
	tokens, err := security.LoadOrGenerate(os.Getenv("JWT_PRIVATE_KEY_PEM"), os.Getenv("JWT_KID"))
	if err != nil {
		log.Fatal(err)
	}

	mux := httpadapter.NewMux(httpadapter.Deps{
		Service: service,
		Auth: app.Auth{
			Users:     users,
			Passwords: security.Passwords{},
			Tokens:    tokens,
		},
	})

	log.Printf("drp-identity-api listening on %s (%s; JWKS issued here)", addr, store)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

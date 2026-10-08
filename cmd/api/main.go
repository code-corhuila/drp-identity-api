package main

import (
	"log"
	"net/http"
	"os"

	httpadapter "github.com/code-corhuila/drp-identity-api/internal/adapters/http"
	"github.com/code-corhuila/drp-identity-api/internal/adapters/memory"
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

	users, err := memory.Corte2(os.Getenv("IDENTITY_SEED_PASSWORD"))
	if err != nil {
		log.Fatal(err)
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

	log.Printf("drp-identity-api listening on %s (memory Corte 2 seed; front does not require this process)", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

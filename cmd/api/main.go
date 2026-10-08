package main

import (
	"log"
	"net/http"
	"os"

	httpadapter "github.com/code-corhuila/drp-identity-api/internal/adapters/http"
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

	log.Printf("drp-identity-api listening on %s", addr)
	if err := http.ListenAndServe(addr, httpadapter.NewMux(service)); err != nil {
		log.Fatal(err)
	}
}

package httpadapter

import (
	"encoding/json"
	"net/http"

	"github.com/code-corhuila/drp-identity-api/internal/app"
)

type healthBody struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}

func HealthHandler(service string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h := app.CheckHealth(service)
		writeJSON(w, http.StatusOK, healthBody{Status: h.Status, Service: h.Service})
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

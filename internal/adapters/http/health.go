package httpadapter

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/code-corhuila/drp-identity-api/internal/app"
)

type healthBody struct {
	Status    string `json:"status"`
	Service   string `json:"service"`
	Timestamp string `json:"timestamp"`
}

func HealthHandler(service string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h := app.CheckHealth(service, time.Now().UTC())
		writeJSON(w, http.StatusOK, healthBody{
			Status:    h.Status,
			Service:   h.Service,
			Timestamp: h.Timestamp.UTC().Format("2006-01-02T15:04:05Z"),
		})
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

package app

import "time"

// Health is the liveness result (E-04: status, service, timestamp).
type Health struct {
	Status    string
	Service   string
	Timestamp time.Time
}

func CheckHealth(service string, now time.Time) Health {
	if service == "" {
		service = "identity-service"
	}
	if now.IsZero() {
		now = time.Now().UTC()
	}
	return Health{Status: "ok", Service: service, Timestamp: now.UTC()}
}

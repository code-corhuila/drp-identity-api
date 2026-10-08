package app

// Health is the liveness result. No database: login persistence lands in a later PR.
type Health struct {
	Status  string
	Service string
}

func CheckHealth(service string) Health {
	if service == "" {
		service = "identity-service"
	}
	return Health{Status: "ok", Service: service}
}

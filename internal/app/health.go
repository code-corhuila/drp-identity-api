package app

// Health is the liveness result. Login uses an in-memory seed until Flyway persistence lands.
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

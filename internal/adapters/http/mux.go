package httpadapter

import "net/http"

func NewMux(service string) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("GET /health", HealthHandler(service))
	mux.Handle("/", NotFoundHandler())
	return WithCorrelation(mux)
}

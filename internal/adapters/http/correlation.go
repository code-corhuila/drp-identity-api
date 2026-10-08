package httpadapter

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
)

type ctxKey string

const correlationKey ctxKey = "traceId"

func WithCorrelation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Correlation-Id")
		if id == "" {
			id = newTraceID()
		}
		w.Header().Set("X-Correlation-Id", id)
		ctx := context.WithValue(r.Context(), correlationKey, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func TraceID(r *http.Request) string {
	if v, ok := r.Context().Value(correlationKey).(string); ok {
		return v
	}
	return ""
}

func newTraceID() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	return hex.EncodeToString(b[:])
}

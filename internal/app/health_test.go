package app

import (
	"testing"
	"time"
)

func TestCheckHealth(t *testing.T) {
	now := time.Date(2026, 11, 3, 9, 0, 0, 0, time.UTC)
	got := CheckHealth("identity-service", now)
	if got.Status != "ok" || got.Service != "identity-service" || !got.Timestamp.Equal(now) {
		t.Fatalf("got %+v", got)
	}
}

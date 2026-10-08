package app

import "testing"

func TestCheckHealth(t *testing.T) {
	got := CheckHealth("identity-service")
	if got.Status != "ok" || got.Service != "identity-service" {
		t.Fatalf("got %+v", got)
	}
}

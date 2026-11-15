package httpadapter

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoginIgnoresUserField(t *testing.T) {
	srv := httptest.NewServer(authMux(t))
	defer srv.Close()

	payload := []byte(`{"email":"member@spacehub.local","password":"Spacehub1!","user":{"id":"spoof"}}`)
	res, err := http.Post(srv.URL+"/api/v1/auth/login", "application/json", bytes.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("login %d", res.StatusCode)
	}
	var raw map[string]any
	if err := json.NewDecoder(res.Body).Decode(&raw); err != nil {
		t.Fatal(err)
	}
	if _, ok := raw["user"]; ok {
		t.Fatalf("login response must not include user: %+v", raw)
	}
	if raw["accessToken"] == "" || raw["tokenType"] != "Bearer" {
		t.Fatalf("%+v", raw)
	}
}

func TestLogin422Envelope(t *testing.T) {
	srv := httptest.NewServer(authMux(t))
	defer srv.Close()

	req, _ := http.NewRequest(http.MethodPost, srv.URL+"/api/v1/auth/login", bytes.NewReader([]byte(`{"email":"member@spacehub.local","password":"wrongpass"}`)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Correlation-Id", "22222222-2222-2222-2222-222222222222")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Fatalf("status %d", res.StatusCode)
	}
	var raw map[string]any
	if err := json.NewDecoder(res.Body).Decode(&raw); err != nil {
		t.Fatal(err)
	}
	if raw["error"] != "BUSINESS_RULE_VIOLATION" || raw["message"] == "" || raw["traceId"] != "22222222-2222-2222-2222-222222222222" {
		t.Fatalf("envelope %+v", raw)
	}
	if _, ok := raw["details"]; ok && raw["details"] == nil {
		t.Fatalf("details present but null: %+v", raw)
	}
}

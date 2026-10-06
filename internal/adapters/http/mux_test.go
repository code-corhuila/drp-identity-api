package httpadapter

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealth(t *testing.T) {
	srv := httptest.NewServer(NewMux("identity-service"))
	defer srv.Close()

	res, err := http.Get(srv.URL + "/health")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status %d", res.StatusCode)
	}
	if res.Header.Get("X-Correlation-Id") == "" {
		t.Fatal("missing X-Correlation-Id")
	}
	var body healthBody
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if body.Status != "ok" || body.Service != "identity-service" {
		t.Fatalf("body %+v", body)
	}
}

func TestNotFoundEnvelope(t *testing.T) {
	srv := httptest.NewServer(NewMux("identity-service"))
	defer srv.Close()

	req, _ := http.NewRequest(http.MethodGet, srv.URL+"/nope", nil)
	req.Header.Set("X-Correlation-Id", "11111111-1111-1111-1111-111111111111")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusNotFound {
		t.Fatalf("status %d", res.StatusCode)
	}
	var body errorBody
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if body.Error != "NOT_FOUND" || body.TraceID != "11111111-1111-1111-1111-111111111111" {
		t.Fatalf("body %+v", body)
	}
}

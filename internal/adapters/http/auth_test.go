package httpadapter

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/code-corhuila/drp-identity-api/internal/adapters/memory"
	"github.com/code-corhuila/drp-identity-api/internal/adapters/security"
	"github.com/code-corhuila/drp-identity-api/internal/app"
)

func authMux(t *testing.T) http.Handler {
	t.Helper()
	users, err := memory.Corte2("Spacehub1!")
	if err != nil {
		t.Fatal(err)
	}
	keys, err := security.LoadOrGenerate("", security.DefaultKID)
	if err != nil {
		t.Fatal(err)
	}
	return NewMux(Deps{
		Service: "identity-service",
		Auth:    app.Auth{Users: users, Passwords: security.Passwords{}, Tokens: keys},
	})
}

func TestLoginAndMe(t *testing.T) {
	srv := httptest.NewServer(authMux(t))
	defer srv.Close()

	body, _ := json.Marshal(map[string]string{"email": "member@spacehub.local", "password": "Spacehub1!"})
	res, err := http.Post(srv.URL+"/api/v1/auth/login", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("login %d", res.StatusCode)
	}
	var got loginResponse
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}
	if got.TokenType != "Bearer" || got.User.Role != "USER" || got.AccessToken == "" {
		t.Fatalf("%+v", got)
	}

	req, _ := http.NewRequest(http.MethodGet, srv.URL+"/api/v1/users/me", nil)
	req.Header.Set("Authorization", "Bearer "+got.AccessToken)
	me, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer me.Body.Close()
	if me.StatusCode != http.StatusOK {
		t.Fatalf("me %d", me.StatusCode)
	}
}

func TestLoginUnauthorizedAndJWKS(t *testing.T) {
	srv := httptest.NewServer(authMux(t))
	defer srv.Close()

	body, _ := json.Marshal(map[string]string{"email": "member@spacehub.local", "password": "wrongpass"})
	res, err := http.Post(srv.URL+"/api/v1/auth/login", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("status %d", res.StatusCode)
	}

	jwks, err := http.Get(srv.URL + "/api/v1/auth/jwks")
	if err != nil {
		t.Fatal(err)
	}
	defer jwks.Body.Close()
	if jwks.StatusCode != http.StatusOK {
		t.Fatalf("jwks %d", jwks.StatusCode)
	}
	var doc map[string]any
	if err := json.NewDecoder(jwks.Body).Decode(&doc); err != nil {
		t.Fatal(err)
	}
	keys, _ := doc["keys"].([]any)
	if len(keys) != 1 {
		t.Fatalf("keys %+v", doc)
	}
}

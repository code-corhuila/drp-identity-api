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

const testSeedPassword = "test-seed-only"

func authMux(t *testing.T) http.Handler {
	t.Helper()
	users, err := memory.Corte2(testSeedPassword)
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

	body, _ := json.Marshal(map[string]string{"email": "member@spacehub.local", "password": testSeedPassword})
	res, err := http.Post(srv.URL+"/api/v1/auth/login", "application/json", bytes.NewReader(body))
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
		t.Fatalf("E-01 login must not include user: %+v", raw)
	}
	token, _ := raw["accessToken"].(string)
	tokenType, _ := raw["tokenType"].(string)
	expires, _ := raw["expiresIn"].(float64)
	if token == "" || tokenType != "Bearer" || expires != 3600 {
		t.Fatalf("%+v", raw)
	}

	req, _ := http.NewRequest(http.MethodGet, srv.URL+"/api/v1/users/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	me, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer me.Body.Close()
	if me.StatusCode != http.StatusOK {
		t.Fatalf("me %d", me.StatusCode)
	}
	var profile userDTO
	if err := json.NewDecoder(me.Body).Decode(&profile); err != nil {
		t.Fatal(err)
	}
	if profile.Role != "USER" || profile.Email != "member@spacehub.local" {
		t.Fatalf("%+v", profile)
	}
}

func TestMeIsAlwaysTheBearerSubject(t *testing.T) {
	srv := httptest.NewServer(authMux(t))
	defer srv.Close()

	memberLogin := postLogin(t, srv.URL, `{"email":"member@spacehub.local","password":"`+testSeedPassword+`"}`)
	adminLogin := postLogin(t, srv.URL, `{"email":"admin@spacehub.local","password":"`+testSeedPassword+`"}`)
	if memberLogin.status != http.StatusOK || adminLogin.status != http.StatusOK {
		t.Fatalf("login member=%d admin=%d", memberLogin.status, adminLogin.status)
	}
	memberTok, _ := memberLogin.body["accessToken"].(string)
	adminTok, _ := adminLogin.body["accessToken"].(string)

	memberMe := getMe(t, srv.URL, memberTok)
	adminMe := getMe(t, srv.URL, adminTok)
	if memberMe.Email != "member@spacehub.local" || memberMe.Role != "USER" {
		t.Fatalf("member me %+v", memberMe)
	}
	if adminMe.Email != "admin@spacehub.local" || adminMe.Role != "ADMIN" {
		t.Fatalf("admin me %+v", adminMe)
	}
	if memberMe.ID == adminMe.ID {
		t.Fatal("two users must not share /users/me")
	}
}

func getMe(t *testing.T, base, token string) userDTO {
	t.Helper()
	req, _ := http.NewRequest(http.MethodGet, base+"/api/v1/users/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("me %d", res.StatusCode)
	}
	var profile userDTO
	if err := json.NewDecoder(res.Body).Decode(&profile); err != nil {
		t.Fatal(err)
	}
	return profile
}

func TestLoginUnprocessableAndJWKS(t *testing.T) {
	srv := httptest.NewServer(authMux(t))
	defer srv.Close()

	wrong := postLogin(t, srv.URL, `{"email":"member@spacehub.local","password":"wrongpass"}`)
	unknown := postLogin(t, srv.URL, `{"email":"nobody@spacehub.local","password":"`+testSeedPassword+`"}`)
	if wrong.status != http.StatusUnprocessableEntity || unknown.status != http.StatusUnprocessableEntity {
		t.Fatalf("status wrong=%d unknown=%d", wrong.status, unknown.status)
	}
	if wrong.body["error"] != "BUSINESS_RULE_VIOLATION" || wrong.body["message"] != unknown.body["message"] {
		t.Fatalf("messages must match: %+v %+v", wrong.body, unknown.body)
	}
	if wrong.body["message"] != "Usuario o contraseña incorrectos." {
		t.Fatalf("message %+v", wrong.body)
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

func TestLoginEmptyVsMissingPassword(t *testing.T) {
	srv := httptest.NewServer(authMux(t))
	defer srv.Close()

	empty := postLogin(t, srv.URL, `{"email":"member@spacehub.local","password":""}`)
	if empty.status != http.StatusUnprocessableEntity {
		t.Fatalf("empty password %d %+v", empty.status, empty.body)
	}
	missing := postLogin(t, srv.URL, `{"email":"member@spacehub.local"}`)
	if missing.status != http.StatusBadRequest {
		t.Fatalf("missing password %d %+v", missing.status, missing.body)
	}
	details, _ := missing.body["details"].([]any)
	if len(details) == 0 {
		t.Fatalf("missing password details %+v", missing.body)
	}
	row, _ := details[0].(map[string]any)
	if row["field"] != "password" {
		t.Fatalf("details field %+v", details)
	}
	nullPwd := postLogin(t, srv.URL, `{"email":"member@spacehub.local","password":null}`)
	if nullPwd.status != http.StatusBadRequest {
		t.Fatalf("null password %d %+v", nullPwd.status, nullPwd.body)
	}
}

type loginHTTP struct {
	status int
	body   map[string]any
}

func postLogin(t *testing.T, base, payload string) loginHTTP {
	t.Helper()
	res, err := http.Post(base+"/api/v1/auth/login", "application/json", bytes.NewReader([]byte(payload)))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	var body map[string]any
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	return loginHTTP{status: res.StatusCode, body: body}
}

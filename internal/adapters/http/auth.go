package httpadapter

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/code-corhuila/drp-identity-api/internal/app"
	"github.com/code-corhuila/drp-identity-api/internal/domain"
)

type loginRequest struct {
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

type userDTO struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	DisplayName string `json:"displayName,omitempty"`
	Role        string `json:"role"`
}

type loginResponse struct {
	AccessToken string `json:"accessToken"`
	TokenType   string `json:"tokenType"`
	ExpiresIn   int    `json:"expiresIn"`
}

func toUserDTO(u domain.User) userDTO {
	return userDTO{ID: u.ID, Email: u.Email.String(), DisplayName: u.DisplayName, Role: string(u.PublicRole())}
}

func LoginHandler(auth app.Auth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req loginRequest
		dec := json.NewDecoder(io.LimitReader(r.Body, 1<<20))
		if err := dec.Decode(&req); err != nil {
			writeErr(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "Datos de entrada inválidos", nil)
			return
		}
		if req.Email == nil || req.Password == nil {
			writeErr(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "Datos de entrada inválidos", []map[string]string{{"field": "email", "message": "Email y contraseña son requeridos"}})
			return
		}
		got, err := auth.Login(r.Context(), app.LoginCommand{Email: *req.Email, Password: *req.Password})
		if errors.Is(err, app.ErrValidation) {
			writeErr(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "Datos de entrada inválidos", []map[string]string{{"field": "email", "message": "Email y contraseña (mín. 8) son requeridos"}})
			return
		}
		if errors.Is(err, app.ErrInvalidCredentials) {
			writeErr(w, r, http.StatusUnprocessableEntity, "BUSINESS_RULE_VIOLATION", "Usuario o contraseña incorrectos.", nil)
			return
		}
		if err != nil {
			writeErr(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "Error interno", nil)
			return
		}
		writeJSON(w, http.StatusOK, loginResponse{
			AccessToken: got.AccessToken,
			TokenType:   got.TokenType,
			ExpiresIn:   got.ExpiresIn,
		})
	}
}

func MeHandler(auth app.Auth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, err := auth.CurrentUser(r.Context(), r.Header.Get("Authorization"))
		if err != nil {
			writeErr(w, r, http.StatusUnauthorized, "UNAUTHORIZED", "Token de autenticación requerido", nil)
			return
		}
		writeJSON(w, http.StatusOK, toUserDTO(u))
	}
}

func JWKSHandler(tokens app.TokenService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, tokens.PublicJWKS())
	}
}

func writeErr(w http.ResponseWriter, r *http.Request, status int, code, msg string, details []map[string]string) {
	body := map[string]any{"error": code, "message": msg, "traceId": TraceID(r)}
	if details != nil {
		body["details"] = details
	}
	writeJSON(w, status, body)
}

func stripSpoofedIdentity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := r.Header.Clone()
		for k := range h {
			if strings.HasPrefix(strings.ToLower(k), "x-user") {
				h.Del(k)
			}
		}
		r.Header = h
		next.ServeHTTP(w, r)
	})
}

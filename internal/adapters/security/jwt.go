package security

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"math/big"
	"os"
	"time"

	"github.com/code-corhuila/drp-identity-api/internal/domain"
	"github.com/golang-jwt/jwt/v5"
)

const DefaultKID = "spacehub-identity-2026"
const DefaultTTL = 3600

type RS256 struct {
	key *rsa.PrivateKey
	kid string
	ttl time.Duration
}

func LoadOrGenerate(pemPath, kid string) (*RS256, error) {
	if kid == "" {
		kid = DefaultKID
	}
	if pemPath != "" {
		raw, err := os.ReadFile(pemPath)
		if err != nil {
			return nil, err
		}
		block, _ := pem.Decode(raw)
		if block == nil {
			return nil, errors.New("invalid pem")
		}
		key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			k, err2 := x509.ParsePKCS8PrivateKey(block.Bytes)
			if err2 != nil {
				return nil, err
			}
			var ok bool
			key, ok = k.(*rsa.PrivateKey)
			if !ok {
				return nil, errors.New("not rsa")
			}
		}
		return &RS256{key: key, kid: kid, ttl: DefaultTTL * time.Second}, nil
	}
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	return &RS256{key: key, kid: kid, ttl: DefaultTTL * time.Second}, nil
}

func (s *RS256) Issue(user domain.User) (string, int, error) {
	now := time.Now().UTC()
	claims := jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email.String(),
		"roles": roles(user),
		"iat":   now.Unix(),
		"exp":   now.Add(s.ttl).Unix(),
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tok.Header["kid"] = s.kid
	signed, err := tok.SignedString(s.key)
	return signed, int(s.ttl.Seconds()), err
}

func (s *RS256) Parse(raw string) (string, error) {
	parser := jwt.NewParser(jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Alg()}))
	tok, err := parser.Parse(raw, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodRS256 {
			return nil, jwt.ErrTokenSignatureInvalid
		}
		return &s.key.PublicKey, nil
	})
	if err != nil || !tok.Valid {
		return "", err
	}
	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok {
		return "", jwt.ErrTokenInvalidClaims
	}
	sub, _ := claims.GetSubject()
	if sub == "" {
		return "", jwt.ErrTokenInvalidClaims
	}
	return sub, nil
}

func (s *RS256) PublicJWKS() map[string]any {
	pub := &s.key.PublicKey
	return map[string]any{
		"keys": []map[string]string{{
			"kty": "RSA",
			"use": "sig",
			"kid": s.kid,
			"alg": "RS256",
			"n":   base64.RawURLEncoding.EncodeToString(pub.N.Bytes()),
			"e":   base64.RawURLEncoding.EncodeToString(big.NewInt(int64(pub.E)).Bytes()),
		}},
	}
}

func roles(u domain.User) []string {
	out := make([]string, 0, len(u.Roles))
	for _, r := range u.Roles {
		out = append(out, string(r))
	}
	return out
}

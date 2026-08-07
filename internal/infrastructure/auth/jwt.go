package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"dach-admin/internal/domain"

	"github.com/google/uuid"
)

type JWT struct {
	secret []byte
	ttl    time.Duration
}

func NewJWT(secret string, ttl time.Duration) *JWT {
	return &JWT{secret: []byte(secret), ttl: ttl}
}

func (j *JWT) Issue(member domain.TeamMember) (string, error) {
	header := map[string]string{"alg": "HS256", "typ": "JWT"}
	claims := domain.AuthClaims{Subject: member.ID.String(), Email: member.Email, Role: member.Role, Exp: time.Now().Add(j.ttl).Unix()}
	return j.sign(header, claims)
}

func (j *JWT) Verify(token string) (domain.AuthClaims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return domain.AuthClaims{}, domain.ErrUnauthenticated
	}
	unsigned := parts[0] + "." + parts[1]
	expected := signature(unsigned, j.secret)
	if !hmac.Equal([]byte(expected), []byte(parts[2])) {
		return domain.AuthClaims{}, domain.ErrUnauthenticated
	}
	var claims domain.AuthClaims
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return domain.AuthClaims{}, domain.ErrUnauthenticated
	}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return domain.AuthClaims{}, domain.ErrUnauthenticated
	}
	if claims.Exp < time.Now().Unix() {
		return domain.AuthClaims{}, domain.ErrUnauthenticated
	}
	if _, err := uuid.Parse(claims.Subject); err != nil {
		return domain.AuthClaims{}, domain.ErrUnauthenticated
	}
	return claims, nil
}

func (j *JWT) sign(header map[string]string, claims domain.AuthClaims) (string, error) {
	h, err := json.Marshal(header)
	if err != nil {
		return "", err
	}
	p, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}
	if len(j.secret) == 0 {
		return "", errors.New("jwt secret is empty")
	}
	unsigned := base64.RawURLEncoding.EncodeToString(h) + "." + base64.RawURLEncoding.EncodeToString(p)
	return unsigned + "." + signature(unsigned, j.secret), nil
}

func signature(unsigned string, secret []byte) string {
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(unsigned))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

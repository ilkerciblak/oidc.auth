package adapter

import (
	"auth-app/internal/service/auth"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type jWTTokenManager struct {
	secret string
}

func JWTTokenManager(secret string) *jWTTokenManager {
	if strings.EqualFold(strings.TrimSpace(secret), "") {
		panic("jwt secret is undefined")
	}

	return &jWTTokenManager{
		secret: secret,
	}
}

func (m *jWTTokenManager) GenerateToken(uid, email, role, issuer, audience string, duration time.Duration) (string, error) {
	claims, err := auth.GetClaims(
		uid,
		email,
		role,
		issuer,
		audience,
		duration,
	)
	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims, nil)
	token_string, err := token.SignedString([]byte(m.secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %v", err)
	}

	return token_string, nil
}

func (m *jWTTokenManager) ValidateToken(token string) (*auth.JWTClaims, error) {
	to, err := jwt.ParseWithClaims(
		token,
		&auth.JWTClaims{},
		func(t *jwt.Token) (any, error) {
			if _, k := t.Method.(*jwt.SigningMethodHMAC); !k {
				return nil, fmt.Errorf("invalid signing method")
			}
		
			return []byte(m.secret), nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %v", err)
	}
	claims, k := to.Claims.(*auth.JWTClaims)
	if !k {
		return nil, fmt.Errorf("unexpected claims format")
	}
	if time.Now().After(claims.ExpiresAt.Time) {
		return nil, fmt.Errorf("token expired")
	}

	return claims, nil
}

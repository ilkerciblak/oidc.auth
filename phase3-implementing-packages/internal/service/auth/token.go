package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenManager interface {
	GenerateToken(uid, email, role, issuer, audience string, duration time.Duration) (string, error)
	ValidateToken(token string) (*JWTClaims, error)
}

type JWTClaims struct {
	UID   string `json:"user_id"`
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}

func GetClaims(uid, email, role, issuer, audience string, duration time.Duration) (*JWTClaims, error) {
	token_id, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("tokenid generation failed: %v", err)
	}

	now := time.Now()

	return &JWTClaims{
		UID:   uid,
		Email: email,
		Role:  role,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        token_id.String(),
			Issuer:    issuer,
			Audience:  jwt.ClaimStrings{audience},
			Subject:   uid,
			ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}, nil
}


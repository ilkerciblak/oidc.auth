package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	UID   string `json:"user_id"`
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}

func GetClaims(uid, email, role string, duration time.Duration) (*Claims, error) {
	token_id, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("[Token ID Generation Error]: %v", err)
	}
	now := time.Now()
	return &Claims{
		UID:   uid,
		Email: email,
		Role:  role,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        token_id.String(),
			Issuer:    "This and that",
			Audience:  jwt.ClaimStrings{"client_id"},
			Subject:   uid,
			ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}, nil
}

type JwtManager struct {
	Secret string
}


func (j *JwtManager) GenerateToken(uid, email, role string, duration time.Duration) ( string,  error) {
	claims, err := GetClaims(
		uid,
		email,
		role,
		duration,
	)
	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)
	token_str, err := token.SignedString([]byte(j.Secret))
	if err != nil {
		return "", err
	}

	return token_str, nil
}

func (m *JwtManager) ValidateToken(token_str string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		token_str,
		&Claims{},
		func(t *jwt.Token) (any, error) {
			_, k := t.Method.(*jwt.SigningMethodRSA)
			if !k {
				return nil, fmt.Errorf("invalid signing method")
			}

			return []byte(m.Secret), nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("Invalid Token: %v", err)
	}

	claims, k := token.Claims.(*Claims)
	if !k {
		return nil, fmt.Errorf("[ValidateToken]: Invalid Claims Format")
	}

	if time.Now().Unix() > claims.ExpiresAt.Unix() {
		return nil, fmt.Errorf("[Token Expired]")
	}

	return claims, nil
}

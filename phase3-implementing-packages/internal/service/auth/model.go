package auth

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
)

type Provider interface {
	GetAuthUrl(state, nonce string) (string, error)
	VerifyIdToken(ctx context.Context, id_token string) (ProviderClaims, error)
	ExchangeCode(ctx context.Context, access_code string) (AuthToken, error)
	GetUserInfo(ctx context.Context, access_token string) (UserInfo, error)
	GetName() string
	DoesSupportOIDC() bool
}

type UserInfo interface {
	GetProviderUID() string
//	Extra(key string) (string, error)
}

type AuthToken interface {
	GetIdToken() string
	GetAccessToken() string
	GetExpiresIn() int
	GetScope() string
	GetTokenType() string
	GetRefreshToken() string
}
type jwk struct {
	Use string `json:"use"`
	Kid string `json:"kid"`
	Alg string `json:"alg"`
	Kty string `json:"kty"`
	E   string `json:"e"`
	N   string `json:"n"`
}

type jwksResponse struct {
	Keys []jwk `json:"keys"`
}

type ProviderClaims interface {
	GetNonce() string
	GetEmail() string
	GetEmailVerified() bool
	GetHd() string
	GetAtHash() string
	GetAuthorizedParty() string
	jwt.Claims
}

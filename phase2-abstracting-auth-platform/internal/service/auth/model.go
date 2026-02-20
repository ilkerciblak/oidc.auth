package auth

import "github.com/golang-jwt/jwt/v5"

type Provider interface {
	GetAuthUrl(state, nonce string) string
	VerifyIdToken(id_token string) string
	ExchangeCode(access_code string) (AuthToken, error)
}

type AuthToken interface {
	IdToken() string
	AccessCode() string
	ExpiresIn() int
	Scope() string
	TokenType() string
	RefreshToken() string
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

type GoogleClaims struct {
	Nonce string `json:"nonce"`
	// The user's email address. Provided only `email` scope is included
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	// The domain associated with the Google Workspace or Cloud organization of the user.
	Hd string `json:"hd"`
	// Access Token Hash. This claim can be used to protect agains xss attacks.
	AtHash string `json:"at_hash"`
	// Identifies the client that the token was issued to (client_id)
	AuthorizedParty string `json:"azp"`

	jwt.RegisteredClaims
}

func (g GoogleClaims) GetNonce() string {
	return g.Nonce
}

func (g GoogleClaims) GetEmail() string {
	return g.Email
}

func (g GoogleClaims) GetEmailVerified() bool {
	return g.EmailVerified
}

func (g GoogleClaims) GetHd() string {
	return g.Hd
}

func (g GoogleClaims) GetAtHash() string {
	return g.AtHash
}

func (g GoogleClaims) GetAuthorizedParty() string {
	return g.AuthorizedParty
}


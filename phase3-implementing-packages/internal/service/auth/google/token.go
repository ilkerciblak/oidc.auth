package google

import (
	"github.com/golang-jwt/jwt/v5"
)

type googleAuthToken struct {
	IdToken      string `json:"id_token"`
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
}

func (g googleAuthToken) GetIdToken() string {
	return g.IdToken
}

func (g googleAuthToken) GetAccessToken() string {
	return g.AccessToken
}

func (g googleAuthToken) GetExpiresIn() int {
	return g.ExpiresIn
}

func (g googleAuthToken) GetScope() string {
	return g.Scope
}

func (g googleAuthToken) GetTokenType() string {
	return g.TokenType
}

func (g googleAuthToken) GetRefreshToken() string {
	return g.RefreshToken
}

type googleAuthClaims struct {
	Nonce           string `json:"nonce,omitempty"`
	Email           string `json:"email,omitempty"`
	EmailVerified   bool   `json:"email_verified,omitempty"`
	Hd              string `json:"hd,omitempty"`
	AtHash          string `json:"at_hash,omitempty"`
	AuthorizedParty string `json:"azp,omitempty"`
	jwt.RegisteredClaims
}

func (g googleAuthClaims) GetNonce() string {
	return g.Nonce
}

func (g googleAuthClaims) GetEmail() string {
	return g.Email
}

func (g googleAuthClaims) GetEmailVerified() bool {
	return g.EmailVerified
}

func (g googleAuthClaims) GetHd() string {
	return g.Hd
}

func (g googleAuthClaims) GetAtHash() string {
	return g.AtHash
}

func (g googleAuthClaims) GetAuthorizedParty() string {
	return g.AuthorizedParty
}

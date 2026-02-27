package facebook

import (
	"github.com/golang-jwt/jwt/v5"
)

type facebookAuthToken struct {
	IdToken      string `json:"id_token"`
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
}

func (g facebookAuthToken) GetIdToken() string {
	return g.IdToken
}

func (g facebookAuthToken) GetAccessToken() string {
	return g.AccessToken
}

func (g facebookAuthToken) GetExpiresIn() int {
	return g.ExpiresIn
}

func (g facebookAuthToken) GetScope() string {
	return g.Scope
}

func (g facebookAuthToken) GetTokenType() string {
	return g.TokenType
}

func (g facebookAuthToken) GetRefreshToken() string {
	return g.RefreshToken
}

type facebookAuthClaims struct {
	Nonce           string `json:"nonce,omitempty"`
	Email           string `json:"email,omitempty"`
	EmailVerified   bool   `json:"email_verified,omitempty"`
	Hd              string `json:"hd,omitempty"`
	AtHash          string `json:"at_hash,omitempty"`
	AuthorizedParty string `json:"azp,omitempty"`
	jwt.RegisteredClaims
}

func (g facebookAuthClaims) GetNonce() string {
	return g.Nonce
}

func (g facebookAuthClaims) GetEmail() string {
	return g.Email
}

func (g facebookAuthClaims) GetEmailVerified() bool {
	return g.EmailVerified
}

func (g facebookAuthClaims) GetHd() string {
	return g.Hd
}

func (g facebookAuthClaims) GetAtHash() string {
	return g.AtHash
}

func (g facebookAuthClaims) GetAuthorizedParty() string {
	return g.AuthorizedParty
}

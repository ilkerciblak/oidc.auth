package github

import "github.com/golang-jwt/jwt/v5"

type githubAuthToken struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
}
func (g githubAuthToken) GetIdToken() string {
	return ""
}

func (g githubAuthToken) GetAccessToken() string {
	return g.AccessToken
}

func (g githubAuthToken) GetExpiresIn() int {
	return g.ExpiresIn
}

func (g githubAuthToken) GetScope() string {
	return g.Scope
}

func (g githubAuthToken) GetTokenType() string {
	return g.TokenType
}

func (g githubAuthToken) GetRefreshToken() string {
	return g.RefreshToken
}

type githubAuthClaims struct {
	Nonce           string `json:"nonce,omitempty"`
	Email           string `json:"email,omitempty"`
	EmailVerified   bool   `json:"email_verified,omitempty"`
	Hd              string `json:"hd,omitempty"`
	AtHash          string `json:"at_hash,omitempty"`
	AuthorizedParty string `json:"azp,omitempty"`
	jwt.RegisteredClaims
}

func (g githubAuthClaims) GetNonce() string {
	return g.Nonce
}

func (g githubAuthClaims) GetEmail() string {
	return g.Email
}

func (g githubAuthClaims) GetEmailVerified() bool {
	return g.EmailVerified
}

func (g githubAuthClaims) GetHd() string {
	return g.Hd
}

func (g githubAuthClaims) GetAtHash() string {
	return g.AtHash
}

func (g githubAuthClaims) GetAuthorizedParty() string {
	return g.AuthorizedParty
}




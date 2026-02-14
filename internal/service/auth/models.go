package auth

type TokenResponse struct {
	IdToken      string `json:"id_token"`
	AccessToken  string `json:"access_code"`
	ExpiresIn    string `json:"expires_in"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
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

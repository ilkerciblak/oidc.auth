package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type googleProvider struct {
	ClientID     string
	ClientSecret string
	AuthURI      string
	TokenURI     string
	RedirectURI  string
	Scopes       []string
}

func (p googleProvider) CreateAuthUrl(state, nonce string) string {
	query_params := url.Values{}
	query_params.Add("client_id", p.ClientID)
	query_params.Add("redirect_uri", p.RedirectURI)
	query_params.Add("scope", strings.Join(p.Scopes, "%20"))
	query_params.Add("response_type", "scope")
	query_params.Add("state", state)
	query_params.Add("nonce", nonce)

	return fmt.Sprintf("%s?%s", p.AuthURI, query_params.Encode())
}

// Please go see [Google OIDC Documentation](developers.google.com/identity/openid-connect/openid-connect) for more claim fields
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

func (p googleProvider) VerifyIdToken(id_token_str string) (*GoogleClaims, error) {
	token, err := jwt.ParseWithClaims(
		id_token_str,
		&GoogleClaims{},
		func(t *jwt.Token) (any, error) {
			val, k := t.Method.(*jwt.SigningMethodRSA)
			if !k || val.Name != jwt.SigningMethodRS256.Name {
				return nil, fmt.Errorf("Invalid Token Signing Method")
			}
			kid, ok := t.Header["kid"].(string)
			if !ok{
				return nil, fmt.Errorf("kid not found in the token header")
			}

			publicKey, k := p.jwtksCache[kid]
			if k! {
				return nil, fmt.Errorf("Public key cannot found for kid=%s", kid)

			return []byte(p.ClientSecret), nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("Invalid Token Format")
	}

	claims, k := token.Claims.(*GoogleClaims)
	if !k {
		return nil, fmt.Errorf("Invalid Claims Format")
	}

	return claims, nil
}

func (g googleProvider) ExchangeCode(access_code_str string) (*TokenResponse, error) {
	params := url.Values{}
	params.Add("code", access_code_str)
	params.Add("client_id", g.ClientID)
	params.Add("cleint_secret", g.ClientSecret)
	params.Add("redirect_uri", g.RedirectURI)
	params.Add("grant_type", "authorization_code")
	r, err := http.NewRequest(
		http.MethodPost,
		g.TokenURI,
		strings.NewReader(params.Encode()),
	)
	if err != nil {
		return nil, err
	}

	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := http.Client{}

	resp, err := client.Do(r)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var tokenResponse TokenResponse

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&tokenResponse); err != nil {
		return nil, err
	}

	return &tokenResponse, nil
}

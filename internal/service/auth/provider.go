package auth

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type googleProvider struct {
	ClientID         string
	ClientSecret     string
	AuthURI          string
	TokenURI         string
	RedirectURI      string
	JWKsURI          string
	Scopes           []string
	cachedPublicKeys map[string]*rsa.PublicKey
	jwtksExpirety    time.Time
}

func (p googleProvider) CreateAuthUrl(state, nonce string) string {
	query_params := url.Values{}
	query_params.Add("client_id", p.ClientID)
	query_params.Add("redirect_uri", p.RedirectURI)
	query_params.Add("scope", strings.Join(p.Scopes, "%20"))
	query_params.Add("response_type", "code")
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

func jwkToRSAPublicKey(key jwk) (*rsa.PublicKey, error) {
	nBytes, err := base64.RawURLEncoding.DecodeString(key.N)
	if err != nil {
		return nil, err
	}
	eBytes, err := base64.RawURLEncoding.DecodeString(key.E)
	if err != nil {
		return nil, err
	}
	n := new(big.Int).SetBytes(nBytes)
	e := new(big.Int).SetBytes(eBytes)
		

	return &rsa.PublicKey{
		N: n,
		E: int(e.Int64()),
	}, nil
}

func (p *googleProvider) fetchJWKS() error {
	if time.Now().Before(p.jwtksExpirety) && p.cachedPublicKeys != nil {
		return nil
	}

	resp, err := http.Get(p.JWKsURI)
	if err != nil {
		return fmt.Errorf("[FAILED TO FETCH JWKS KEYS]: %v", err)
	}
	defer resp.Body.Close()

	var jwks jwksResponse
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return fmt.Errorf("[FAILED TO DECODE JWKs RESPONSE]: %v", err)
	}

	p.cachedPublicKeys = make(map[string]*rsa.PublicKey)
	for _, key := range jwks.Keys {
		pubKey, err := jwkToRSAPublicKey(key)
		if err != nil {
			return err
		}

		p.cachedPublicKeys[key.Kid] = pubKey
	}

	p.jwtksExpirety = time.Now().Add(1 * time.Hour)

	return nil
}

func (p googleProvider) VerifyIdToken(id_token_str string) (*GoogleClaims, error) {
	// Fetch JWKS Keys
	if err := p.fetchJWKS(); err != nil {
		return nil, err
	}

	token, err := jwt.ParseWithClaims(
		id_token_str,
		&GoogleClaims{},
		func(t *jwt.Token) (any, error) {
			val, k := t.Method.(*jwt.SigningMethodRSA)
			if !k || val.Name != jwt.SigningMethodRS256.Name {
				return nil, fmt.Errorf("Invalid Token Signing Method")
			}

			kid, k := t.Header["kid"].(string)
			if !k {
				return nil, fmt.Errorf("Invalid Token: kid is not in the header")
			}

			publicKey, exists := p.cachedPublicKeys[kid]
			if !exists {
				return nil, fmt.Errorf("Invalid Token: public key not found for kid (%s)", kid)
			}

			return publicKey, nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("Invalid Token Format")
	}

	claims, k := token.Claims.(*GoogleClaims)
	if !k {
		return nil, fmt.Errorf("Invalid Claims Format")
	}

	if !strings.EqualFold("https://account.google.com", claims.Issuer)|| strings.EqualFold(claims.Issuer, "account.google.com") {
		return nil, fmt.Errorf("Invalid Token: invalid issuer ")
	}

	if claims.Audience[0] != p.ClientID {
		return nil, fmt.Errorf("Invalid Audience")
	}

	if time.Now().Unix() > claims.ExpiresAt.Unix() {
		return nil, fmt.Errorf("Token Expired")
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

	client := http.Client{
		Timeout: time.Duration(10 * time.Second),
	}

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

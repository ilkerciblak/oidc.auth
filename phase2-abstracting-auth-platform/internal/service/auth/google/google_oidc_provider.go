package google

import (
	"auth-app/internal/service/auth"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type googleOIDCProvider struct {
	config *auth.ProviderConfig
}

func GoogleOIDCProvider(configFuncs ...auth.WithProviderConfig) *googleOIDCProvider {
	cfg := &auth.ProviderConfig{}

	for _, f := range configFuncs {
		cfg = f(cfg)
	}

	cfg.Discover()
	cfg.FetchCacheKeys()

	return &googleOIDCProvider{
		config: cfg,
	}
}

func (g *googleOIDCProvider) DoesSupportOIDC() bool {
	return true
}

func (g *googleOIDCProvider) GetName() string {
	return "google"
}

func (p *googleOIDCProvider) GetAuthUrl(state, nonce string) (string, error) {
	query_params := url.Values{}
	query_params.Add("client_id", p.config.ClientID)
	query_params.Add("redirect_uri", p.config.CallbackURI)
	query_params.Add("scope", p.config.GetScopes())
	query_params.Add("response_type", "code")
	query_params.Add("state", state)
	query_params.Add("nonce", nonce)

	return fmt.Sprintf("%s?%s", p.config.AuthURI, query_params.Encode()), nil
}

func (p *googleOIDCProvider) VerifyIdToken(id_token string) (auth.ProviderClaims, error) {
	if err := p.config.FetchCacheKeys(); err != nil {
		return nil, err
	}

	token, err := jwt.ParseWithClaims(
		id_token,
		&googleAuthClaims{},
		func(t *jwt.Token) (any, error) {
			if _, k := t.Method.(*jwt.SigningMethodRSA); !k {
				return nil, fmt.Errorf("invalid signing method")
			}

			kid, k := t.Header["kid"].(string)
			if !k {
				return nil, fmt.Errorf("kid not found in the token header")
			}

			public_key, k := p.config.CachedPublicKeys[kid]
			if !k {
				return nil, fmt.Errorf("public key not found for kid: %s", kid)
			}
			return public_key, nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse claims %v", err)
	}

	claims, k := token.Claims.(*googleAuthClaims)
	if !k {
		return nil, fmt.Errorf("invalid claims format")
	}

	if !strings.EqualFold(claims.Issuer, "https://accounts.google.com") && !strings.EqualFold(claims.Issuer, "accounts.google.com") {
		return nil, fmt.Errorf("invalid issuer")
	}

	if claims.Audience[0] != p.config.ClientID {
		return nil, fmt.Errorf("invalid audience")
	}

	if time.Now().Unix() > claims.ExpiresAt.Unix() {
		return nil, fmt.Errorf("token expired")
	}

	return claims, nil
}

func (p *googleOIDCProvider) ExchangeCode(access_code string) (auth.AuthToken, error) {
	params := url.Values{}
	params.Add("code", access_code)
	params.Add("client_id", p.config.ClientID)
	params.Add("client_secret", p.config.ClientSecret)
	params.Add("redirect_uri", p.config.CallbackURI)
	params.Add("grant_type", "authorization_code")
	req, err := http.NewRequest(
		http.MethodPost,
		p.config.TokenURI,
		strings.NewReader(params.Encode()),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("content-type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	client.Timeout = 10 * time.Second

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("token exchange failed %v", err)
	}
	defer resp.Body.Close()
	var auth_token googleAuthToken
	if err := json.NewDecoder(resp.Body).Decode(&auth_token); err != nil {
		return nil, fmt.Errorf("token decoding failed %v", err)
	}

	return auth_token, nil
}

package facebook

import (
	"auth-app/internal/service/auth"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type facebookAuthProvider struct {
	config *auth.ProviderConfig
}

func FacebookAuthProvider(config_funcs ...auth.WithProviderConfig) *facebookAuthProvider {
	cfg := &auth.ProviderConfig{}

	for _, f := range config_funcs {
		cfg = f(cfg)
	}

	return &facebookAuthProvider{
		config: cfg,
	}
}

func (p *facebookAuthProvider) GetAuthUrl(state, nonce string) (string, error) {
	params := url.Values{}
	params.Add("client_id", p.config.ClientID)
	params.Add("scope", p.config.GetScopes())
	params.Add("redirect_uri", p.config.CallbackURI)
	params.Add("response_type", "code")
	params.Add("state", state)
	params.Add("nonce", nonce)
	/// Code challenge is not required any more in backend confidential clients
	//	code_challenge, err := platform.GenerateBase64EncodedString()
	//	if err != nil {
	//		return "", fmt.Errorf("failed to generate code_challenge: %v", err)
	//	}
	//	params.Add("code_challenge", code_challenge)
	//	params.Add("code_challenge_method", "plain")

	return fmt.Sprintf("%s?%s", p.config.AuthURI, params.Encode()), nil
}

func (p *facebookAuthProvider) VerifyIdToken(id_token string) (auth.ProviderClaims, error) {

	if err := p.config.FetchCacheKeys(); err != nil {
		return nil, err
	}

	token, err := jwt.ParseWithClaims(
		id_token,
		&facebookAuthClaims{},
		func(t *jwt.Token) (any, error) {
			if _, k := t.Method.(*jwt.SigningMethodRSA); !k {
				return nil, fmt.Errorf("invalid signing method")
			}

			kid, k := t.Header["kid"].(string)
			if !k {
				return nil, fmt.Errorf("kid is not found")
			}

			public_key, k := p.config.CachedPublicKeys[kid]
			if !k {
				return nil, fmt.Errorf("public is undefined for kid=%s", kid)
			}

			return public_key, nil

		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token %v", err)
	}

	claims, k := token.Claims.(*facebookAuthClaims)
	if !k {
		return nil, fmt.Errorf("invalid claims format")
	}

	if claims.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("token is expired")
	}

	return claims, nil
}

func (p *facebookAuthProvider) ExchangeCode(access_code string) (auth.AuthToken, error) {

	params := url.Values{}
	params.Add("client_id", p.config.ClientID)
	params.Add("client_secret", p.config.ClientSecret)
	params.Add("redirect_uri", p.config.CallbackURI)
	params.Add("code", access_code)

	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf(
			"%s?%s",
			p.config.TokenURI,
			params.Encode(),
		),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to construct request: %v", err)
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %v", err)
	}

	defer res.Body.Close()

	var tkn facebookAuthToken

	if err := json.NewDecoder(res.Body).Decode(&tkn); err != nil {
		return nil, fmt.Errorf("token decoding failed: %v", err)
	}

	return tkn, nil
}

func (p *facebookAuthProvider) GetName() string {
	return "facebook"
}
func (p *facebookAuthProvider) DoesSupportOIDC() bool {
	return true
}

package github

import (
	"auth-app/internal/service/auth"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type githubOAuthProvider struct {
	config *auth.ProviderConfig
}

func GithubOAuthProvider(configs ...auth.WithProviderConfig) *githubOAuthProvider {
	cfg := &auth.ProviderConfig{}

	for _, f := range configs {
		cfg = f(cfg)
	}

	return &githubOAuthProvider{
		config: cfg,
	}
}

func (p *githubOAuthProvider) DoesSupportOIDC() bool {
	return false
}

func (p *githubOAuthProvider) GetName() string {
	return "github"
}

func (p *githubOAuthProvider) GetAuthUrl(state, nonce string) (string, error) {
	params := url.Values{}
	params.Add("client_id", p.config.ClientID)
	params.Add("redirect_uri", p.config.CallbackURI)
	//params.Add("login", "true")
	params.Add("scope", p.config.GetScopes())
	params.Add("state", state)

	return fmt.Sprintf("%s?%s", p.config.AuthURI, params.Encode()), nil

}

func (p *githubOAuthProvider) VerifyIdToken(id_token string) (auth.ProviderClaims, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		"https://api.github.com/user",
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to construct request:%v", err)
	}

	req.Header.Set(
		"Authorization",
		fmt.Sprintf("Bearer %s", id_token),
	)

	req.Header.Set(
		"Accept",
		"application/vnd.github+json",
	)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	//bodyBytes, _ := io.ReadAll(res.Body)
	//	fmt.Println(string(bodyBytes))

	var claims struct {
		UserID int `json:"id"`
	}

	if err := json.NewDecoder(res.Body).Decode(&claims); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}


	return githubAuthClaims{

		RegisteredClaims: jwt.RegisteredClaims{
			Subject: fmt.Sprintf("%d", claims.UserID),
		},
	}, nil
}

func (p *githubOAuthProvider) ExchangeCode(access_code string) (auth.AuthToken, error) {
	params := url.Values{}
	params.Add("client_id", p.config.ClientID)
	params.Add("client_secret", p.config.ClientSecret)
	params.Add("code", access_code)

	req, err := http.NewRequest(
		http.MethodPost,
		p.config.TokenURI,
		strings.NewReader(params.Encode()),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("content-type", "application/x-www-form-urlencoded")
	req.Header.Set("accept", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("token exchange failed: %v", err)
	}
	defer res.Body.Close()

	var auth_token githubAuthToken

	if err := json.NewDecoder(res.Body).Decode(&auth_token); err != nil {
		return nil, fmt.Errorf("token decoding failed %v", err)
	}

	return auth_token, nil
}

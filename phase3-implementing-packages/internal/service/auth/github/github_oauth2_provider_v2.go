package github

import (
	"auth-app/internal/service/auth"
	"context"
	"fmt"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/endpoints"
)

type githubOauth2Providerv2 struct {
	provider *oidc.Provider
	config   *oauth2.Config
}

func GithubOauth2Providerv2(ctx context.Context, client_id, client_secret, redirect_url string) *githubOauth2Providerv2 {
	provider_cfg := oidc.ProviderConfig{
		IssuerURL:     "https://github.com/login/oauth/",
		AuthURL:       endpoints.GitHub.AuthURL,
		TokenURL:      endpoints.GitHub.TokenURL,
		DeviceAuthURL: "",
		UserInfoURL:   "https://api.github.com/user",
		JWKSURL:       "https://github.com/login/oauth/.well-known/jwks",
		Algorithms:    []string{"RS256"},
	}
	provider := provider_cfg.NewProvider(ctx)

	config := &oauth2.Config{
		ClientID:     client_id,
		ClientSecret: client_secret,
		RedirectURL:  redirect_url,
		Endpoint:     endpoints.GitHub,
		Scopes:       []string{"read:user"},
	}

	return &githubOauth2Providerv2{
		provider: provider,
		config:   config,
	}

}

func (p *githubOauth2Providerv2) GetAuthUrl(state, nonce string) (string, error) {
	return p.config.AuthCodeURL(state), nil
}
func (p *githubOauth2Providerv2) VerifyIdToken(ctx context.Context, id_token string) (auth.ProviderClaims, error) {
	panic("not implemented")
}
func (p *githubOauth2Providerv2) ExchangeCode(ctx context.Context, access_code string) (auth.AuthToken, error) {
	tkn, err := p.config.Exchange(
		ctx,
		access_code,
	)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %v", err)
	}

	token := githubAuthToken{
		AccessToken:  tkn.AccessToken,
		ExpiresIn:    int(tkn.ExpiresIn),
		TokenType:    tkn.TokenType,
		RefreshToken: tkn.RefreshToken,
	}

	return token, nil

}
func (p *githubOauth2Providerv2) GetName() string {
	return "github"
}
func (p *githubOauth2Providerv2) DoesSupportOIDC() bool {
	return false
}
func (p *githubOauth2Providerv2) GetUserInfo(ctx context.Context, access_token string) (auth.UserInfo, error) {
	user_info, err := p.provider.UserInfo(ctx, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: access_token}))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user_info: %v", err)
	}
	var claims GitHubUserInfo
	if err := user_info.Claims(&claims); err != nil {
		return nil, fmt.Errorf("failed to parse user info: %v", err)
	}

	return &claims, nil
}

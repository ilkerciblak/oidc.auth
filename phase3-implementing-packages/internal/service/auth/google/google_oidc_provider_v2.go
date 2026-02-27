package google

import (
	"auth-app/internal/service/auth"
	"context"
	"fmt"

	"github.com/coreos/go-oidc/v3/oidc"

	"golang.org/x/oauth2"
)

type googleOIDCProviderv2 struct {
	provider *oidc.Provider
	verifier *oidc.IDTokenVerifier
	cfg      *oauth2.Config
}

func GoogleOIDCProviderV2(ctx context.Context, clientID, clientSecret, redirect_url string) *googleOIDCProviderv2 {
	provider, err := oidc.NewProvider(ctx, "https://accounts.google.com")
	if err != nil {
		panic(fmt.Errorf("failed to construct google oidc provider: %v", err))
	}

	cfg := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  redirect_url,
		Scopes:       []string{oidc.ScopeOpenID, "email"},
	}

	verifier := provider.Verifier(&oidc.Config{
		ClientID: cfg.ClientID,
	})

	return &googleOIDCProviderv2{
		provider: provider,
		verifier: verifier,
		cfg:      cfg,
	}
}

func (p googleOIDCProviderv2) GetAuthUrl(state, nonce string) (string, error) {

	return p.cfg.AuthCodeURL(state), nil

}
func (p googleOIDCProviderv2) VerifyIdToken(ctx context.Context, id_token string) (auth.ProviderClaims, error) {
	token, err := p.verifier.Verify(
		ctx,
		id_token,
	)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	var claims googleAuthClaims
	if err := token.Claims(&claims); err != nil {
		return nil, fmt.Errorf("claims parse failed: %v", err)
	}

	return claims, nil
}
func (p googleOIDCProviderv2) ExchangeCode(ctx context.Context, access_code string) (auth.AuthToken, error) {
	token, err := p.cfg.Exchange(ctx, access_code, oauth2.AccessTypeOnline)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %v", err)
	}
	tkn := &googleAuthToken{
		IdToken:      token.Extra("id_token").(string),
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresIn:    int(token.ExpiresIn),
		Scope:        token.Extra("scope").(string),
		TokenType:    token.TokenType,
	}

	return tkn, nil

}

func (p googleOIDCProviderv2) GetName() string {
	return "google"
}

func (p googleOIDCProviderv2) DoesSupportOIDC() bool {
	return true
}

func (p *googleOIDCProviderv2) GetUserInfo(ctx context.Context, access_token string) (auth.UserInfo, error) {

	i, err := p.provider.UserInfo(
		ctx,
		oauth2.StaticTokenSource(&oauth2.Token{
			AccessToken: access_token,
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user info: %v", err)
	}
	var claims GoogleUserInfo
	if err := i.Claims(
		&claims,
	); err != nil {
		return nil, fmt.Errorf("failed to parse userinfo : %v", err)
	}

	return &claims, nil
}

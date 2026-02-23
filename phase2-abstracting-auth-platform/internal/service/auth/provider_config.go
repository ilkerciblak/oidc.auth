package auth

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"time"
)

type ProviderConfig struct {
	ClientID         string
	ClientSecret     string
	DiscoverURI      string
	CallbackURI      string
	Scopes           []string
	CachedPublicKeys map[string]*rsa.PublicKey
	KeyExpirety      time.Time
	AuthURI          string
	TokenURI         string
	JWKsURI          string
}

type WithProviderConfig func(*ProviderConfig) *ProviderConfig

func WithClientID(client_id string) WithProviderConfig {
	return func(pc *ProviderConfig) *ProviderConfig {
		pc.ClientID = client_id
		return pc
	}
}

func WithClientSecret(client_id string) WithProviderConfig {
	return func(pc *ProviderConfig) *ProviderConfig {
		pc.ClientSecret = client_id
		return pc
	}
}
func WithDiscoverURI(client_id string) WithProviderConfig {
	return func(pc *ProviderConfig) *ProviderConfig {
		pc.DiscoverURI = client_id
		return pc
	}
}
func WithCallbackURI(client_id string) WithProviderConfig {
	return func(pc *ProviderConfig) *ProviderConfig {
		pc.CallbackURI = client_id
		return pc
	}
}
func WithScopes(scopes []string) WithProviderConfig {
	return func(pc *ProviderConfig) *ProviderConfig {
		pc.Scopes = scopes 
		return pc
	}
}




func (p *ProviderConfig) GetScopes() string {
	return strings.Join(p.Scopes, " ")
}

func (p *ProviderConfig) Discover() {
	if strings.EqualFold(strings.TrimSpace(p.DiscoverURI), "") {
		panic(fmt.Errorf("DiscoverURI is empty"))
	}
	var resp struct {
		AuthURI  string `json:"authorization_endpoint"`
		TokenURI string `json:"token_endpoint"`
		JwksURI  string `json:"jwks_uri"`
	}

	res, err := http.Get(p.DiscoverURI)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		panic(fmt.Errorf("decoding discovery:%v", err))
	}

	if resp.AuthURI == "" || resp.TokenURI == "" || resp.JwksURI == "" {
		panic(fmt.Sprintf("discover document is not sufficient for %s, set URIs manually", p.DiscoverURI))
	}
	p.AuthURI = resp.AuthURI
	p.TokenURI = resp.TokenURI
	p.JWKsURI = resp.JwksURI
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

func (p *ProviderConfig) FetchCacheKeys() error {
	if p.CachedPublicKeys != nil && time.Now().Before(p.KeyExpirety) {
		return nil
	}

	res, err := http.Get(p.JWKsURI)
	if err != nil {
		return fmt.Errorf("failed to fetch jwks: %v", err)
	}
	defer res.Body.Close()

	var jwks jwksResponse
	if err := json.NewDecoder(res.Body).Decode(&jwks); err != nil {
		return fmt.Errorf("failed to decode jwks response %v", err)
	}
	p.CachedPublicKeys = make(map[string]*rsa.PublicKey)
	for _, key := range jwks.Keys {
		pubKey, err := jwkToRSAPublicKey(key)
		if err != nil {
			return err
		}
		p.CachedPublicKeys[key.Kid] = pubKey
	}

	if cacheControl := res.Header.Get("cache-control"); cacheControl != "" {
		var max_age int
		_, err := fmt.Sscanf(cacheControl, "max-age=%d", &max_age)
		if err != nil || max_age <= 0 {
			p.KeyExpirety = time.Now().Add(time.Hour * 1)
		}
		p.KeyExpirety = time.Now().Add(time.Duration(max_age) * time.Second)
	}
	return nil
}

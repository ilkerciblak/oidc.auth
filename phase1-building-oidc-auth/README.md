

# Phase1: Building Simple Authentication Server with OIDC Support

## Overview

In *Phase 1*, development survey will include building simple authentication server with `Google OIDCsupport` implemented. Abstraction of the oidc flow will be skipped to present whole mechanism. Thus development will be architected as only `GoogleOIDCProvider` will be implemented and used. On the other hand, SOLID principles are covered at most.

Following parts of this documentation, reader can figure out how the `oidc` flows around the project.Documentation will split the presentation to layers of project which are 

 
## Project Structure: Phase 1.0.0
```bash
phase1-building-oidc-auth
├── go.mod
├── go.sum
├── internal
│   ├── platform
│   │   └── config.go # Environment Variables Manager
│   └── service
│       ├── auth
│       │   ├── handler.go # LoginWGoogle, Callback endpoints
│       │   ├── models.go # Domain Entites
│       │   ├── provider.go # Google Provider Client Model
│       │   ├── state_manager.go # State Manager 
│       │   └── token.go # Token Manager
│       └── user
├── main.go # Main entrence of the application
└── README.md
```

## Features

- **Google OIDC Login**: Secure user login with Google Support

## Tech Stack and Libraries
In phase 1.0.0, project only involves of a **HTTP RESTful Server API** developed with *Golang* with following packages:
- `net/http` - HTTP server/client
- `encoding/json` - JSON parse
- `encoding/base64` - Base64 encode/decode
- `crypto/rsa` - RSA signature verification
- `crypto/rand` - Random string generation
- `crypto/sha256` - Hashing
- `time` - Expiration check

## Required Environment Variables
```.env
# APP
HOST=required
PORT=required

# PROVIDERS
## GOOGLE OIDC
GOOGLE_CLIENT_ID=required
GOOGLE_CLIENT_SECRET=required
GOOGLE_PROJECT_ID=required
GOOGLE_X509_CERT_URL=optional_for_local_dev
GOOGLE_AUTH_URI=required
GOOGLE_TOKEN_URI=required
GOOGLE_REDIRECT_URI=required
GOOGLE_JWKS_URI=required

JWT_SECRET=required
```
See [Google Cloud Client Page]() to create and/or obtain your client information. _cannot give mine 🦆._

## Architectural Overview

Since authentication flow described in the [main README file](../README.md). This documentation aimed to explain the package structure.

### Presentation Layer: GoogleOIDCHandler 

In order to adopt *dependency inversion*, `GoogleOIDCHandler` is the object where package deals it dependency injections. Also client facing endpoints, such as `Login and Callback` will be defined here.

GoogleOIDCHandler will have three dependencies, 

- The GoogleProvider object
- StateManager to create and validate `state` parameter
- JWTManager to create and validate `JWT`

In addition to that, handler functions
- **Login**: The handler receives end-users sign-in demand, generating `nonce` and `state` parameters and redirecting client to the oidc provider's auth screen.
- **Callback**: The `redirect_url` that oidc providers mandated to request with `authorization code` and `state` parameters. Within this handler function, `state` parameter can be validated and `authorization code exchanging process` can be done. After code exchange process we can decode the `id_token`and provide `access_token` and `refresh_token` to our end-user with redirecting them to some `home` page.


```go
type GoogleOIDCHandler struct {
	GoogleProvider
	*StateManager
	*JwtManager
}

func (g GoogleOIDCHandler) Login(w http.ResponseWriter, r *http.Request) {
	state, err := g.GenerateState()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	nonce := ""
	authUrl := g.CreateAuthUrl(
		state,
		nonce,
	)

	http.Redirect(
		w,
		r,
		authUrl,
		http.StatusFound,
	)
}

func (g GoogleOIDCHandler) Callback(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	// nonce := r.URL.Query().Get("nonce")
	code := r.URL.Query().Get("code")

	if err := g.ValidateState(state); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
	}

	g_token, err := g.ExchangeCode(code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
	}

	g_claims, err := g.VerifyIdToken(g_token.IdToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	token_str, err := g.GenerateToken(
		g_claims.Subject,
		g_claims.Email,
		g_claims.Email,
		time.Duration(time.Minute*5),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	refresh_token, err := g.GenerateToken(
		g_claims.Subject,
		g_claims.Email,
		g_claims.Email,
		time.Duration(time.Minute*5),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	http.SetCookie(
		w,
		&http.Cookie{
			Name:     "access_token",
			Value:    token_str,
			SameSite: http.SameSiteStrictMode,
			Secure:   false, // development environment helloo 🙋🏽‍♂️
			HttpOnly: true,
		},
	)

	http.SetCookie(
		w,
		&http.Cookie{
			Name:     "refresh_token",
			Value:    refresh_token,
			SameSite: http.SameSiteStrictMode,
			Secure:   false, // development environment helloo 🙋🏽‍♂️
			HttpOnly: true,
		},
	)

	http.Redirect(w, r, "/home", http.StatusFound)
}
```

### Application Layer: GoogleProvider

For Phase 1.x, abstractions are skipped. Thus application layer will only consists of `GoogleProvider` instrumentation. Since provider discovery mechanism is not used yet, endpoints are hard-coded to `.env` file. Consequently, using environment variables, provider attributes will be set.

For now, GoogleProvider model will provide `CreateAuthUrl`, `VerifyIdToken` and `ExchangeCode` publicmethods, which are constructs oidc authentication flow for that provider.

One important point to explain, as  [OpenID Connection Specs](https://openid.net/specs/openid-connect-core-1_0.html#CodeIDToken) declares, an `id_token` exchanged for `authorization_code` is a JWT token that provides standard and provider specific (based on `scope` parameter) user claims to usein client party to authenticate end-user. Client party must validate and decrypt the `id_token` with the signature in the algorithm specified in the provider JWS response.   

```go
type TokenResponse struct {
	IdToken      string `json:"id_token"`
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
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
```


In [Google OIDC Document](https://developers.google.com/identity/openid-connect/openid-connect) it ismentioned that _Google changes public keys only **in-frequently**_. Thus these public keys are can be cached. Thus following instrument have `jwksExpirety` attribute to set TTL for cached public keys.


```go
type GoogleProvider struct {
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

func (p GoogleProvider) CreateAuthUrl(state, nonce string) string {
	query_params := url.Values{}
	query_params.Add("client_id", p.ClientID)
	query_params.Add("redirect_uri", p.RedirectURI)
	query_params.Add("scope", strings.Join(p.Scopes, " "))
	query_params.Add("response_type", "code")
	query_params.Add("state", state)w
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

func (p *GoogleProvider) fetchJWKS() error {
	if p.cachedPublicKeys != nil && time.Now().Before(p.jwtksExpirety) {
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
	if cacheControl := resp.Header.Get("Cache-Control"); cacheControl != "" {
		var max_age int
		_, err := fmt.Sscanf(cacheControl, "max-age=%d", &max_age)
		if err != nil || max_age <= 0 {
			p.jwtksExpirety = time.Now().Add(1 * time.Hour)
		}

		p.jwtksExpirety = time.Now().Add(time.Duration(max_age) * time.Second)

	}

	return nil
}

func (p GoogleProvider) VerifyIdToken(id_token_str string) (*GoogleClaims, error) {
	// Fetch JWKS Keys
	if err := p.fetchJWKS(); err != nil {
		return nil, err
	}

	token, err := jwt.ParseWithClaims(
		id_token_str,
		&GoogleClaims{},
		func(t *jwt.Token) (any, error) {
			_, k := t.Method.(*jwt.SigningMethodRSA)
			if !k {
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

	//if !strings.EqualFold("https://account.google.com", claims.Issuer) || strings.EqualFold(claims.Issuer, "account.google.com") {
	//	return nil, fmt.Errorf("Invalid Token: invalid issuer ")
	//}

	if claims.Audience[0] != p.ClientID {
		return nil, fmt.Errorf("Invalid Audience")
	}

	if time.Now().Unix() > claims.ExpiresAt.Unix() {
		return nil, fmt.Errorf("Token Expired")
	}

	return claims, nil
}

func (g GoogleProvider) ExchangeCode(access_code_str string) (*TokenResponse, error) {
	params := url.Values{}
	params.Add("code", access_code_str)
	params.Add("client_id", g.ClientID)
	params.Add("client_secret", g.ClientSecret)
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
```

## Versioning & Changelog

In order to implement _changes_, e.g. refactors or new features, in several steps this repository will use `Phase x.x.x` declaration. 

Thus to follow steps and further changes, please check [issues](https://github.com/ilkerciblak/oidc.auth/issues) and [CHANGELOG.md](../CHANGELOG.md) 

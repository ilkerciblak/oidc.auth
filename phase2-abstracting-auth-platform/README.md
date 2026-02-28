

# Phase2: Implementing Abstraction Layers to Authentication Service

## Overview

In order to introduce testable, design and implementation efficiency, and encapsulation, in *Phase 2* some abstraction will be implemented to our authenticatin service. 

Since `oidc and oauth2` protocols introduces same rules to follow up, the authentication business domain can be split into some interfaces and domain types to use with multiple implementations. According to that, authentication flow can be split into following concepts

- **Provider Model Interface**: Encapsulates the authorization endpoint construction, code exchange and id_token validation processes by means of each provider suggestions

- **Provider Config Struct**: Provides re-usable configuration mechanism for all provides to configure variables as `client_id and secrets` also `jwks public keys`
- **State Store Interface**: Encapsulates the `state parameter` generation, validation and deletion processes.
- **User Manager Interface**: Encapsulates the `user business logic`. Since this project only scopes user authentication with federated social providers, our interface will only consist of `find or create new user` logic.
- **Token Manager Interface**: Encapsulates the `token based authentication` business logic. In order to provide secure authentication, after federated login logic our application will return its own token to user.

Furthermore, in `adapter/` reader can findout the interface implementations for *State Manager and Token Manager*. Provider implementations will be placed under their own folders e.g. `google/` or `github/`.

 
## Project Structure: Phase 2.0.0
```bash
├── go.mod
├── go.sum
├── internal
│   ├── platform
│   │   ├── config.go
│   │   ├── database.go
│   │   └── str_ops.go
│   └── service
│       ├── auth
│       │   ├── adapter
│       │   │   ├── jwt_token_manager.go
│       │   │   └── redis_state_store.go
│       │   ├── facebook
│       │   │   ├── facebook_oidc_provider.go
│       │   │   └── token.go
│       │   ├── github
│       │   │   ├── github_oauth2_provider.go
│       │   │   └── token.go
│       │   ├── google
│       │   │   ├── google_oidc_provider.go
│       │   │   └── token.go
│       │   ├── handler.go
│       │   ├── model.go
│       │   ├── provider_config.go
│       │   ├── state_store.go
│       │   ├── token.go
│       │   └── user_manager.go
│       └── user
│           ├── model.go
│           └── repository.go
├── main.go
├── migrations
│   └── 00001_init_user.sql
├── README.md
```

## Features

- **Google OIDC Login**: Secure user login with Google Support
- **GitHub OAuth2 Login using UserInfo Endpoint: 

## Tech Stack and Libraries
In phase 1.0.0, project only involves of a **HTTP RESTful Server API** developed with *Golang* with following packages:
- `net/http` - HTTP server/client
- `encoding/json` - JSON parse
- `encoding/base64` - Base64 encode/decode
- `crypto/rsa` - RSA signature verification
- `crypto/rand` - Random string generation
- `crypto/sha256` - Hashing
- `time` - Expiration check

## QuickStart

### Prerequisities

- **Docker**: for development environment and running in container 

### Installation and Running
1. Clone the repository
```bash
git clone https://github.com/ilkerciblak/oidc.auth
```

2. Set up the environment variables

- Create `.env` file in project directory
```bash
cd project-directory && touch .env
```

### Required Environment Variables
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

3. Environment setup and running the container

```bash
cd project-directory
docker compose -f dev.docker-compose.yml up -d
```

4. Running the phase2 using Makefile in container
```bash
cd project-directory
docker exec -it auth-dev make run-2
```

## Architectural Overview

### How it works?

In order to provide follow along documentation, an example implementation will be given.

### Implementing Google Sign In Feature with OIDC 

1. Provider interface implies
```go
type Provider interface {
    // Construct the provider specific authentication_url or returns error
	GetAuthUrl(state, nonce string) (string, error)
    // Encapsulates provider specific id_token verification logic
	VerifyIdToken(id_token string) (ProviderClaims, error)
    // Encapsulates code exchange business logic
	ExchangeCode(access_code string) (AuthToken, error)
    // returns provider name in order to use in user business logic
	GetName() string
    // returns boolean information about whether provider is oidc or oauth2 based 
	DoesSupportOIDC() bool
}
```

2. To implement this interface, our `GoogleOIDCProvider` looks like [this](./internal/service/auth/google/google_oidc_provider.go)

Our `GoogleOIDCProvider` has only one dependency in `ProviderConfig` type. `ProviderConfig` struct introduces an important public method that provides `fetching and caching provider public keys` besides public setter methods. These setter methods will be used in provider definition part.

3. Since `oidc` declares a set of rules for authentication flow, instead of defining a `handler interface` decided to define a `handler struct` with `Login and Callback` re-usable methods. See [handler struct](./internal/service/auth/handler.go)

That approach provided re-usable handler functions with clear dependency injections for each provider. On the other hand, only downside of this approach was occuring with `oauth2 based providers`. Since OAuth2 based providers does not issue `id_token` in result pf code exchange process, their `user info endpoint` should be used to retrieve user information details. ([See the ADR about it](../phase3-imlementing-packages/docs/adr/ADR-0002-Provider-Interface-Change.md)) 

4. Provider instrumentation in `main.go`
```go
// dependencies
state_manager := adapter.RedisStateStore("redis_State:6379", "", 0)
jwt_manager := adapter.JWTTokenManager(cfg.JWT_SECRET)

// provider
google_provider := google.GoogleOIDCProvider(
	auth.WithClientID(cfg.GOOGLE_CLIENT_ID),
	auth.WithClientSecret(cfg.GOOGLE_CLIENT_SECRET),
	auth.WithCallbackURI(cfg.GOOGLE_REDIRECT_URI),
	auth.WithScopes([]string{"openid", "email"}),
	auth.WithDiscoverURI(cfg.GOOGLE_DISCOVER_URI),
)

// handler
google_oidc_handler := auth.OIDCHandler{
	StateManager: state_manager,
	Provider:     google_provider,
	TokenManager: jwt_manager,
	UserManager:  &userRepo,
}


//endpoints
http.HandleFunc("GET /google/auth", google_oidc_handler.Login)
http.HandleFunc(
	`/auth/google/callback`,
	google_oidc_handler.Callback,
)

```

## Versioning & Changelog

In order to implement _changes_, e.g. refactors or new features, in several steps this repository will use `Phase x.x.x` declaration. 

Thus to follow steps and further changes, please check [issues](https://github.com/ilkerciblak/oidc.auth/issues) and [CHANGELOG.md](../CHANGELOG.md) 

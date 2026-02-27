# Phase 3: Building Simple Authentication Server with OIDC Support


## Phase 2.0.0 to 3.0.0 Overview:

In *Phase 3.0.0*, [oidc](https://pkg.go.dev/github.com/coreos/go-oidc/v3/oidc) and [oauth2](https://pkg.go.dev/golang.org/x/oauth2)  packages are implemented.

These migrations improved application's maintability and standards compliance. The `go-oidc` package introduces a more robust `Provider` abstraction that automatically discovers and caches `oidc configurations` through provider' `.well-known` endpoint. Also, `go-oidc` provides provider specific `id token verification` with built-in `JWKS` fetching and rotation support that eliminating manual validation logics.

## Project Structure
```bash
phase3-implementing-packages
├── docs
│   └── adr
│       ├── ADR-0001-Custom-Provider-Token-Usage.md
│       └── ADR-0002-Provider-Interface-Change.md
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
│       │   ├── github
│       │   │   ├── github_oauth2_provider_v2.go
│       │   │   ├── token.go
│       │   │   └── user-info.go
│       │   ├── google
│       │   │   ├── google_oidc_provider_v2.go #provider implementations
│       │   │   ├── token.go # provider-token wrapper
│       │   │   └── user_info.go # user-info-wrapper
│       │   ├── handler.go
│       │   ├── model.go
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

- **Login with Google Account**
- **Login with GitHub Account** 

## Required Environment Variables
```.env
PORT=required

# OIDC and OAuth Apps
## GOOGLE
GOOGLE_CLIENT_ID=required
GOOGLE_CLIENT_SECRET=required
GOOGLE_PROJECT_ID=required
GOOGLE_REDIRECT_URI=required
	
## Github
GITHUB_CLIENT_ID=required
GITHUB_CLIENT_SECRET=required

## Facebook
FACEBOOK_CLIENT_SECRET=required
FACEBOOK_CLIENT_ID=required
FACEBOOK_REDIRECT_URI=required
FACEBOOK_DISCOVER_URI=required

JWT_SECRET=required
	
DB_URL=required

```
## Changelog
 - please check [issues](https://github.com/ilkerciblak/oidc.auth/issues) and [CHANGELOG.md](../CHANGELOG.md) 

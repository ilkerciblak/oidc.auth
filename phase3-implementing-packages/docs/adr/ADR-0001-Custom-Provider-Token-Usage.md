# ADR-0001: Custom ProviderToken Interface for OAuth2 Token Abstraction

**Status:** Accepted  
**Date:** 2026-02-25  
**Deciders:** @ilkerciblak

## Context

The authentication server manages OAuth2-based authentication flows. Go's standard
`golang.org/x/oauth2` package provides an `oauth2.Token` struct, but it does not
fully meet our application's requirements.

## Decision

Instead of using `oauth2.Token` directly, we define our own `ProviderToken` interface:

```go
    type AuthToken interface {
	    GetIdToken() string # not presents in oauth2.Token
	    GetAccessToken() string
    	GetExpiresIn() int
    	GetScope() string
    	GetTokenType() string
    	GetRefreshToken() string
    }
```

## Rationale

Two primary reasons drove this decision:

**1. Abstraction and flexibility**  
Depending directly on `oauth2.Token` would require changing all token-handling code
whenever the provider changes (Google → GitHub → custom OIDC). With this interface,
each provider supplies its own implementation and the core logic remains untouched.

**2. Missing id_token field**  
`oauth2.Token` does not carry the OIDC `id_token` value. This token is critical for
verifying user identity (sub, email, claims). The package exposes it only through
`Extra("id_token")`, which is type-unsafe. Our interface enforces it as a strongly-typed,
required field.

## Alternatives Considered

| Alternative | Reason Rejected |
|---|---|
| Use `oauth2.Token` directly | No `id_token`, provider lock-in risk |
| Embed and extend `oauth2.Token` | Still coupled to a concrete type; harder to mock |
| Use `oidc.IDToken` | Token management and identity split across two structs |

## Consequences

**Positive:**
- Provider implementations are isolated from each other
- Easily mockable in tests
- `id_token` accessible in a type-safe manner

**Negative / Trade-offs:**
- Every new provider requires a wrapper implementation
- Future convenience methods added to `oauth2.Token` are not directly available

## Related Decisions
 - None

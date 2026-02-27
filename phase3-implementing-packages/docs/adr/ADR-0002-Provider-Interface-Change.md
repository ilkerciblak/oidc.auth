# ADR-0002: Provider Interface Changes with UserInfo Method Instrumentation

**Status:** Experimental
**Date:** 2026-02-27 
**Deciders:** @ilkerciblak 


## Context

Our user management flow relies on `provider_name` and `provider_uid`  to check whether a user has previously authenticated with a given identity source. However, certain OAuth2-based platforms — GitHub being a notable example — do not fully conform to the OIDC protocol and lack authorization code flow support, meaning no `id_token` carrying the required attributes is issued. In such cases, UserInfo Endpoints serve as an alternative to obtain the necessary data.

## Decision

Adding `GetUserInfo` method to our `Provider` interface:
```go
    type Provider interface {
    	GetAuthUrl(state, nonce string) (string, error)
    	VerifyIdToken(ctx context.Context, id_token string) (ProviderClaims, error)
    	ExchangeCode(ctx context.Context, access_code string) (AuthToken, error)
    	GetUserInfo(ctx context.Context, access_token string) (UserInfo, error)
    	GetName() string
    	DoesSupportOIDC() bool
    }
```

Instead of using `oidc.UserInfo` directly, we define our own `UserInfo` interface:
```go
    type UserInfo interface {
    	GetProviderUID() string
    }
```

## Rationale

Primary reason drove this decision is:

**Lack of id_token with OAuth2.0 based platforms**:
.
OAuth2.0-based platforms do not fully conform to the authorization code flow, meaning the tokens they issue do not carry an `id_token` containing the required user information. However, their  `UserInfo endpoints` can be used as an alternative to retrieve the necessary details.

## Alternative Considered
- Using ValidateIDToken interface method with different behavior for `oauth2` based platforms

While repurposing the ValidateIDToken interface method was an option, it would have introduced ambiguity into the authentication flow for OAuth2-based platforms, as the method name implies ID token validation — a concept that does not apply in an OAuth2-only context. Additionally, this approach would have hindered the ongoing refactor of our handler's Callback methods rather than complementing it. Beyond resolving this immediate requirement, fetching user information from both OIDC and OAuth2-based platforms adds broader value to most applications, even when user profile data falls outside the current application scope.

## Consequences

**Positive:** 
- Provider implementations are isolated from each other
- Easly mockable in tests
- Method usages are very clear

**Trade-offs:**
- Every new provider requires wrapper implementations

## Related Decisions
- [Custom Token Interface Usage](./ADDR-0001-Custom-Provider-Token-Usage.md)


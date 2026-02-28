# Backend Learning Path: OAuth2 and OIDC Protocols
> [!NOTE]
> Scope of this repository to **Learn while implementing** the `OpenID Connection Protocol`. 
> Thus this documentation will contain a summarized information about `OIDC Authentication Protocol`and 3 phased implementation for `authentication using the authorization code flow`
> 1. Phase 1: Creating an authentication server for Google Sign In with `Golang`, without using `oidc and ouath2` package
>   -  Phase 1.1 : Implementing redis to handle some storing issues instead of using mutexes
>   -  Phase 1.2 : Implementing OpenID Connect Discovery 1.0
>   -  Phase 1.3 : Implementing user business logic
> 2. Phase 2: Refactoring the application with some abstraction layers and implementing additional oidc provider e.g. `GitHub`
> 3. Phase 3: Refactoring the application with implementing `oauth2` and `oidc` package


> [!WARNING]
> 
> This project is not `production ready` due no extra security layers are included e.g. no observability or proper exception handling implemented. This repository only includes the R&D on `OpenID Connect Protocol`. 

## Overview

As mentioned, scope of this project is to *learning while building* the `OpenID Connection and OAuth2.0` topics. Following that concept, this documentation will only cover [How OIDC Protocol Works](## OIDC: OpenID Connect Protocol). In order to follow implementation details reader can visit each implementation documentations which are:
    - [Phase1 implementation details](./phase1-building-oidc-auth/README.md) 
    - [Phase2 implementation details](./phase2-abstracting-auth-platform/README.md) 
    - [Phase3 implementation details](./phase3-implementing-packages/README.md)

Also, differences are issued in [CHANGLELOG](./CHANGELOG.md), to provide clear version history.

## OIDC: OpenID Connect Protocol
OpenID Connection Protocol, OIDC, is an identity layer wrapped aroundthe `OAuth 2.0` framework. While OAuth2 mainly aiming providing resource authorization, main purpose of the OIDC is to providing authentication over federated platforms. Once user demands to log with federated identity, they got re-directed to that third parties auth page to identify themselves and return back with their basic profile information. Furthermore, besides authentication and basic profile information different scopes can be used to reach further user data resource.

### How OpenID Connect Protocol Works - Architectural Structure

OpenID Connect authentication can be performed one of three paths, 
    - `the Authorization Code Flow` where `response_type=code`
    - `the Implicit Flow` where `response_type=id_token%20token` or `response_type=id_token`
    - `the Hybrid Flow`
Mainly the choice of `response_type` parameter is determing the flow type. Flow type determines how the `id_token` and `access_token` are returned to the client.

*[On OpenID specs page](https://openid.net/specs/openid-connect-core-1_0.html#Authenticaiton)* the characteristics of the three flows are summarized as follows. 
| Property | Authorization Code Flow | Implicit Flow | Hybrid Flow |
| --------------- | --------------- | --------------- | --------------- |
| All Tokens returned from **Authorization Endpoint** | No | Yes | No |
| All tokens returned from **Token Endpoint**  | Yes | No | No |
| Tokens are **exposed** to **User Agent**  | No | Yes | Yes |
| Client can be authenticated | Yes | No | Yes |
| Refresh Token is possible | Yes | No | Yes |
| Most communication is **server-to-server** | Yes | No | Varies|

The flow used is determined by the `response_type` value contained in the *Authorization Request*. These `response_type` values select these follows:
| response_type | Flow |
| -------------- | --------------- |
| code | Authorization Code |
| id_token | Implicit Flow |
| id_token token | Implicit Flow |
| code id_token | Hybrid Flow |
| code token | Hybrid Flow |
| code id_token token | Hybrid Flow |

#### Authentication using the Authorization Code Flow

When using the *Authorization Code Flow* all tokens are returned from the *Token Endpoint*. 

The Authorization Code Flow returns an `authorization_code` to the Client, which can then exchange it for an `id_token` and `access_token` directly from *Token Endpoint*. This provides the benefit of *not exposing any tokens to the User Agent.* This will provide security over malicious applications with access to the User Agent.


```mermaid
sequenceDiagram 

	User-->>Frontend:  Demands Sign with Google
	Frontend->>Backend:  GET auth/google
	Note over Backend:	Generate & Redis persists State 
	Backend->>Frontend:  307 Redirect Url
	Frontend-->>User: Browser Redirect to OIDC Page
	participant O as OICD Provider
	User->>O: Login + Allow
	O->>Frontend: 307 Redirect to callback?code=xxx&state=yyy
	Frontend->>Backend: GET /callback?code=xxx&state=yyy
	Note over Backend: Validate state
	Backend->>O: Code - Token Exchange
	O->>Backend: returns id_token+access_token
	Note over Backend: Verify id_token
	Note over Backend: Find or Create user
	Note over Backend: JWT token generate
	Note over Backend: Set Cookie
	Backend->>Frontend: 308 Redirect to home page or any
	Frontend-->>User: Navigate to home page

```

Given mermaid diagram describes the architectural flow for *OIDC Authorization Code Flow*. Each time an `user` demand to `sign-in` with an `OIDC` provider, application’s backend will redirect them to their desired OpenID site where they login with their federated credentials. 

If visitor successfully login and allow this application to use their given scopes, the `OIDC` provider will redirect them with some artifacts. Demanding application will use these given artifacts to validate user session with issuing a request to identity provider again.

If validation process goes well, consequently application can now find an existing user or create a new one with user_info and generate application wise tokens (optional). 

Finally, application can complete authentication process and redirect user to some home page or dashboard etc.

> [!IMPORTANT]
> One important artifact that application will get at the end of an authentication process is `id_token`. This token will contain a set of personal attributes about authenticated user. These attributes can be used in varying use-cases, most importantly identifying the existing user. Most providers will provide `provider_user_id` or `provider_uid` that can be persisted and used to identify whether it belongs to an existing user.

### The ID Token

The primary extension that OIDC makes to OAuth2.0 is to enable End-Users to be **Authenticated** with `id_token` data structure. The `id_token` is a security token that contains `Claims` about the *Authentication of an End-User* by an Authorization Server when using a Client and potentially other requested Claims. The `id_token` is represented as a `JWT`.

- `iss` (Issuer)
Identifies the authentication server that issued the token. Used to verify the token comes from the expected authority.

- `sub` (Subject)
A unique and stable identifier for the user within the issuer. It is used by the client to recognize the user.

- `aud` (Audience)
Specifies which client (client_id) this token is intended for. The client must reject the token if it does not match its own client_id.

- `exp` (Expiration Time)
The time after which the token must no longer be accepted. Prevents the use of old or stolen tokens.

- `iat` (Issued At)
The time when the token was created. Helps determine how old the token is.

- `auth_time`
The time when the user actually authenticated (logged in). Used when enforcing session age (e.g., with max_age).

- `nonce`
A value sent by the client and returned in the token to prevent replay attacks. The client must verify it matches the original request.

- `acr` (Authentication Context Class Reference)
Indicates the assurance or security level of the authentication performed (e.g., MFA vs basic login).

- `amr` (Authentication Methods References)
Lists the authentication methods used (e.g., password, OTP). Shows how the user was authenticated.

- `azp` (Authorized Party)
Identifies the client that the token was issued to, mainly used when multiple audiences are present. Ensures the correct client is authorized to use the token.

ID Tokens are commonly signed using `JWS`, thus in order to reach out these (and other propobal claims) `id_token` should be de-crypted using provided JWS signature. Commonly the signing party publishes its key in a JWK Set at its `jwks_uri` location that includes the `kid` of the signing key. These `public_keys` are rotated *kindly* thus demanding party can use caching for these values. TTL can be found in theprovider's `jwks_uri`s `Cache-Control` header.

### Obtaining User Information from the ID Token

As mentioned, an `id_token` is a `JWT(JWAT)`, where cryptographically signed Base64-encoded JSON object. Normally it is critical to **validate an ID Token** before using it, most oidc providers combine validation with the work of decoding the base64url-encoded values and parsing the JSON within. Thus in most cases, validation process end up with claim access or vice versa.

#### Validating an `id_token`

ID Tokens are sensitive and can be misused if intercepted. Thus these tokens are handled securly by transmitting them only over HTTPS and only using POST data or within the request headers. On the other hand, if your server stores it, it should be stored most securely.

> One thing that makes ID tokens useful is that fact that you can pass them around different components of your app. These components can use an ID token as a lightweight authentication mechanism authenticating the app and the user. But before you can use the information in the ID token or rely on it as an assertion that the user has authenticated, you must validate it.

Validation *- and decryption mostly -* requires several steps:
    1. Verify that the ID Token is properly signed by the issuer. Public Keys can be found using provider specified `jwks_uri`.
    2. Verify `iss` claims with the provider
    3. Verify `aud` claim is equal to your app's client ID
    4. Verify the token expirety, `exp` claim
    5. Verify other provider specific security claims, e.g. `Google` provides `hd` claim to be check.

### Obtaining (more) User Informations

OIDC's provides standart claims, to see please check [The OpenID Connection Specs Page](https://openid.net/specs/openid-connect-core-1_0.html#StandardClaims), which defines a standard set of basic profile Claims. Besides this set, additional claims can be requested using specific `scope` values or individual claims can be requested using the `claims` request parameter. These `scope` values mostly provider specific. For an example, `Google` provides three different scopes in use which are `openid email profile`.

#### UserInfo Endpoint and UserInfo Request

In order to gain access to desired protected resource, the first `auth request` must include related scope. That way your application can use the `access_token` and the `oidc standard`.

After that to gain access to protected resources, just make an HTTPS `GET` request to provider's `userinfo endpoint`. The `userinfo` response includes information about the protected resource as described in `OIDC Standard Claims` and the providers `claims supported` metadata which introduces in their `Discovery Document`. Users or their organizations may choose to supply or withhold certain fields, so you might not get information for every field for your authorized scopes of access.

## Discovery Document, `.well-known`

The OpenID Connect protocol requires the use of multiple endpoints for authenticating users, and for requesting resources including tokens, user information, and public keys.

To simplify implementations and increase flexibility, OpenID Connect allows the use of a "Discovery document," a JSON document found at a well-known location containing key-value pairs which provide details about the OpenID Connect provider's configuration, including the URIs of the authorization, token, revocation, userinfo, and public-keys endpoints. The Discovery document for Google's OpenID Connect service may be retrieved from:

```
https://accounts.google.com/.well-known/openid-configuration
```

example response, which is the google's discovery document
```
{
  "issuer": "https://accounts.google.com",
  "authorization_endpoint": "https://accounts.google.com/o/oauth2/v2/auth",
  "device_authorization_endpoint": "https://oauth2.googleapis.com/device/code",
  "token_endpoint": "https://oauth2.googleapis.com/token",
  "userinfo_endpoint": "https://openidconnect.googleapis.com/v1/userinfo",
  "revocation_endpoint": "https://oauth2.googleapis.com/revoke",
  "jwks_uri": "https://www.googleapis.com/oauth2/v3/certs",
  "response_types_supported": [
    "code",
    "token",
    "id_token",
    "code token",
    "code id_token",
    "token id_token",
    "code token id_token",
    "none"
  ],
  "subject_types_supported": [
    "public"
  ],
  "id_token_signing_alg_values_supported": [
    "RS256"
  ],
  "scopes_supported": [
    "openid",
    "email",
    "profile"
  ],
  "token_endpoint_auth_methods_supported": [
    "client_secret_post",
    "client_secret_basic"
  ],
  "claims_supported": [
    "aud",
    "email",
    "email_verified",
    "exp",
    "family_name",
    "given_name",
    "iat",
    "iss",
    "locale",
    "name",
    "picture",
    "sub"
  ],
  "code_challenge_methods_supported": [
    "plain",
    "S256"
  ]
}
```


## Phase [1.3.0]

### Added
- User business logic implemented, 
    - Only covers *find or create user from provider*. 
- Database connection is instrumented in order to manage user data.

### Changed 
- Provider handler dependecies: `UserRepository` injected.
- Provider handler `Callback` method now checks if the user already signed in with provider and provider_uid.
- Dev. environment: Postgres and Adminer images added.
- Config platform, `DB_URL` attribute is added.


## Phase [1.2.0]

### Added
- OIDC Provider Configuration Discovery 1.0 implemented

### Fixed
- Provider `VerifyIDToken` issuer check

## Phase [1.1.0]

### Added
- Redis store support to state manager
- Static html file for home page, *trivial*  

### Deprecated
- StateManager with in-memory store 

### Refactored
- State generating function re-located under platform, renamed as `GenerateBase64EncodedString`


## Phase [1.0.0] 

### Added
- StateManager instrumented using `sync` mutexes and in memory state store
- OIDC provider instrumented for only Google Accounts
- Provider configurations e.g. endpoints are hard coded in `.env` file

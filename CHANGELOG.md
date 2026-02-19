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

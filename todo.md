## Phase 1:

#### Creating `.env` variables
```markdown
# APP
HOST= fill_
PORT= fill_

# PROVIDERS
## GOOGLE OIDC
GOOGLE_CLIENT_ID= fill_
GOOGLE_CLIENT_SECRET= fill_
GOOGLE_PROJECT_ID= fill_
GOOGLE_X509_CERT_URL= fill_
GOOGLE_AUTH_URI= fill_
GOOGLE_TOKEN_URI= fill_
GOOGLE_REDIRECT_URI= fill_
```


#### Claims Interface
// GetExpirationTime implements the Claims interface.
func (c RegisteredClaims) GetExpirationTime() (*NumericDate, error) {
	return c.ExpiresAt, nil
}

// GetNotBefore implements the Claims interface.
func (c RegisteredClaims) GetNotBefore() (*NumericDate, error) {
	return c.NotBefore, nil
}

// GetIssuedAt implements the Claims interface.
func (c RegisteredClaims) GetIssuedAt() (*NumericDate, error) {
	return c.IssuedAt, nil
}

// GetAudience implements the Claims interface.
func (c RegisteredClaims) GetAudience() (ClaimStrings, error) {
	return c.Audience, nil
}

// GetIssuer implements the Claims interface.
func (c RegisteredClaims) GetIssuer() (string, error) {
	return c.Issuer, nil
}

// GetSubject implements the Claims interface.
func (c RegisteredClaims) GetSubject() (string, error) {
	return c.Subject, nil
}

### Implement some user logic

#### Impelement Redis on phase 2 for state and nonce store

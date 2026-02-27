package google

import (
	"time"
)

type GoogleUserInfo struct {
	Aud           string    `json:"aud,omitempty"`
	Email         string    `json:"email,omitempty"`
	EmailVerified bool      `json:"email_verified,omitempty"`
	Exp           time.Time `json:"exp"`
	FamilyName    string    `json:"family_name,omitempty"`
	GivenName     string    `json:"given_name,omitempty"`
	Iat           string    `json:"iat,omitempty"`
	Iss           string    `json:"iss,omitempty"`
	Locale        string    `json:"locale,omitempty"`
	Name          string    `json:"name,omitempty"`
	Picture       string    `json:"picture,omitempty"`
	Sub           string    `json:"sub,omitempty"`
}

func (p *GoogleUserInfo) GetProviderUID() string {
	return p.Sub
}

/*
Extra(key string) (string, error)
 type UserInfo interface {
}*/

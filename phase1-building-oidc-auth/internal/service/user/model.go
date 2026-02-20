package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id          string    `json:"id"`
	DisplayName string    `json:"display_name"`
	Email       string    `json:"email"`
	Phone       string    `json:"phone"`
	LastLoginAt time.Time `json:"last_login_at"`
	IsVerified  bool      `json:"is_verified"`
}

type AuthProvider struct {
	Id            string `json:"id"`
	Uid           string `json:"uid"`
	Provider      string `json:"provider"`
	ProviderUid   string `json:"provider_uid"`
	LastLoginWith time.Time `json:"last_login_with"`
}

func NewUser(display_name, email, phone string) *User {
	return &User{
		Id:          uuid.New().String(),
		DisplayName: display_name,
		Email:       email,
		Phone:       phone,
		LastLoginAt: time.Now(),
		IsVerified:  true,
	}
}

func NewProvider(uid, provider, provider_uid string) *AuthProvider{

	return &AuthProvider{
		Id: uuid.NewString(),
		Uid: uid,
		Provider: provider,
		ProviderUid: provider_uid,
		LastLoginWith: time.Now(),
	}
}

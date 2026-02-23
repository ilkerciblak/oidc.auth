package auth

import (
	"context"
	"time"
)

type UserManager interface {
	FindOrCreateUser(
		ctx context.Context,
		provider,
		provider_uid string,
	) (User, error)
}

type User interface {
	GetId() string
	GetEmail() string
	GetRole() string
	GetDisplayName() string
	GetPhone() string
	GetLastLoginAt() time.Time 
	GetIsVerified() bool 
	NewUser() User
}

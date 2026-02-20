package auth

import "context"

type UserManager interface {
	FindOrCreateUser(
		ctx context.Context,
		provider,
		provider_uid string,
	) (User, error)
}

type User interface {
	GetId() string
	NewUser() User
}

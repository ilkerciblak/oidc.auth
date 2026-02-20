package auth

import "context"


type StateManager interface{
	Generate(ctx context.Context, provider_name string) string
	Delete(ctx context.Context, state, provider_name string) error
	Validate(ctx context.Context, state, provider_name string) error
}


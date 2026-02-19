package auth

import (
	"auth-app/internal/platform"
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStateManager struct {
	client *redis.Client
}

func NewRedisStateManager() *RedisStateManager {
	client := redis.NewClient(
		&redis.Options{
			Addr:     "redis_state:6379",
			Password: "",
			DB:       0,
		},
	)

	return &RedisStateManager{
		client: client,
	}
}

func (r *RedisStateManager) GenerateState(ctx context.Context, provider_name string) (string, error) {
	state, err := platform.GenerateBase64EncodedString()
	if err != nil {
		return "", fmt.Errorf("[Failed to generate state:] %v", err)
	}

	if err := r.client.Set(
		ctx,
		fmt.Sprintf("%s_oidc:state:%s", provider_name, state),
		state,
		5*time.Minute,
	).Err(); err != nil {
		return "", fmt.Errorf("[Failed to store state:] %v", err)
	}

	return state, nil
}

func (r *RedisStateManager) ValidateState(ctx context.Context, state, provider_name string) error {
	_, err := r.client.Get(
		ctx,
		fmt.Sprintf("%s_oidc:state:%s", provider_name, state),
	).Result()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("state not found")
		}
		return fmt.Errorf("[Failed to Validate State:] %v", err)
	}

	return nil
}

func (r *RedisStateManager) DeleteState(ctx context.Context, state, provider_name string) error {
	if err := r.client.Del(

		ctx,
		fmt.Sprintf("%s_oidc:state:%s", provider_name, state),
	); err != nil {
		return fmt.Errorf("[Failed to Delete State:] %v", err)
	}

	return nil
}

package adapter

import (
	"auth-app/internal/platform"
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisStateStore struct {
	client *redis.Client
}

const key string = "%s_oidc:state:%s"

func RedisStateStore(addr, password string, defaultDb int) *redisStateStore {
	client := redis.NewClient(
		&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       defaultDb,

		},
	)

	return &redisStateStore{
		client: client,
	}
}

func (s *redisStateStore) Generate(ctx context.Context, provider_name string) string {
	state, err := platform.GenerateBase64EncodedString()
	if err != nil {
		return ""
	}

	if err := s.client.Set(
		ctx,
		fmt.Sprintf(key, provider_name, state),
		state,
		5*time.Minute,
	).Err(); err != nil {
		return ""
	}

	fmt.Println("giden state", state)
	return state
}

func (s *redisStateStore) Delete(ctx context.Context, state, provider_name string) error {
	if err := s.client.Del(
		ctx,
		fmt.Sprintf(key, provider_name, state),
	).Err(); err != nil {
		return err
	}

	return nil
}

func (s *redisStateStore) Validate(ctx context.Context, state, provider_name string) error {
	if err := s.client.Get(
		ctx,
		fmt.Sprintf(
			key,
			provider_name,
			state,
		),
	).Err(); err != nil {
		return fmt.Errorf("redis validate failed: %v", err) 
	}
	return nil
}

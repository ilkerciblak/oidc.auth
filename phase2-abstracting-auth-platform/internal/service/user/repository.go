package user

import (
	"auth-app/internal/service/auth"
	"context"
	"database/sql"
	"fmt"
	"time"
)

type UserRepository struct {
	Db *sql.DB
}

var createUserQuery = `
INSERT INTO USERS (id, display_name, email, phone, last_login_at, verified)
VALUES ($1,$2,$3,$4,$5,$6);
`

var findUserQuery = `
SELECT u.*
FROM users AS u
JOIN auth_provider as ap
ON ap.user_id = u.id
WHERE ap.provider=$1 AND ap.provider_user_id=$2;
`

var createAuthProviderQuery = `
INSERT INTO auth_provider (id, user_id, provider, provider_user_id, last_login_with)
VALUES ($1,$2,$3,$4,$5);
`

var updateLastLoginWith = `
UPDATE auth_provider
SET last_login_with=$1
WHERE provider=$2 and provider_user_id=$3;
`

func (r *UserRepository) FindOrCreateUser(ctx context.Context, provider, provider_uid string) (  auth.User, error) {
		fail := func(err error) (  auth.User, error) {
		return nil, err
	}

	// Begin transaction
	tx, err := r.Db.BeginTx(ctx, nil)
	if err != nil {
		return fail(fmt.Errorf("Failed to start transaction: %v", err))
	}
	defer tx.Rollback()
	var user User
	fallbackUser := user.NewUser() 
	// Look for provider and user
	if err := tx.QueryRowContext(
		ctx,
		findUserQuery,
		provider,
		provider_uid,
	).Scan(
		&user.Id,
		&user.DisplayName,
		&user.Email,
		&user.Phone,
		&user.LastLoginAt,
		&user.IsVerified,
	); err != nil {
		if err == sql.ErrNoRows {
			// if no-rows, create user
			if _, err := tx.ExecContext(
				ctx,
				createUserQuery,
				fallbackUser.GetId(),
				fallbackUser.GetDisplayName(),
				fallbackUser.GetEmail(),
				fallbackUser.GetPhone(),
				fallbackUser.GetLastLoginAt(),
				fallbackUser.GetIsVerified(),
			); err != nil {
				return fail(fmt.Errorf("[CreateUser@Tx]: %v", err))
			}

			// create auth-provider
			newProvider := NewProvider(fallbackUser.GetId(), provider, provider_uid)
			if _, err := tx.ExecContext(
				ctx,
				createAuthProviderQuery,
				newProvider.Id,
				fallbackUser.GetId(),
				provider,
				provider_uid,
				time.Now(),
			); err != nil {
				return fail(err)
			}
			if err := tx.Commit(); err != nil {
				return fail(err)
			}
			return  fallbackUser, nil

		}
		return fail(err)
	}

	if _, err := tx.ExecContext(
		ctx,
		updateLastLoginWith,
		time.Now(),
		provider,
		provider_uid,
	); err != nil {
		return fail(err)
	}

	if err := tx.Commit(); err != nil {
		return fail(err)
	}
	return &user, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, user User) (*User, error) {
	_, err := r.Db.Exec(
		createUserQuery,
		user.Id,
		user.DisplayName,
		user.Email,
		user.Phone,
		user.LastLoginAt,
		user.IsVerified,
	)
	if err != nil {
		return nil, fmt.Errorf("[CreateUser]: %v", err)
	}

	return &user, nil
}

func (r *UserRepository) FindUserByProvider(ctx context.Context, provider, provider_uid string) (*User, error) {
	var user User
	if err := r.Db.QueryRow(
		findUserQuery,
		provider,
		provider_uid,
	).Scan(
		&user.Id,
		&user.DisplayName,
		&user.Email,
		&user.Phone,
		&user.LastLoginAt,
		&user.IsVerified,
	); err != nil {
		if err == sql.ErrNoRows {
			return &User{}, nil
		}

		return nil, fmt.Errorf("[FindUserByProvider]: %v", err)
	}

	return &user, nil
}

func (r *UserRepository) CreateAuthProvider(ctx context.Context, provider AuthProvider) (*AuthProvider, error) {
	_, err := r.Db.Exec(
		createAuthProviderQuery,
		provider.Id,
		provider.Uid,
		provider.Provider,
		provider.ProviderUid,
		provider.LastLoginWith,
	)
	if err != nil {
		return nil, fmt.Errorf("[CreateAuthProvider] : %v", err)
	}

	return &provider, nil
}

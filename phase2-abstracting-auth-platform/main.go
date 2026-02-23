package main

import (
	"auth-app/internal/platform"
	"auth-app/internal/service/auth"
	"auth-app/internal/service/auth/adapter"
	"auth-app/internal/service/auth/google"
	"auth-app/internal/service/user"
	"context"
	"fmt"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func main() {
	cfg := platform.LoadConfig()
	db, err := platform.Instrument(
		context.Background(),
		cfg.DB_URL,
	)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := goose.Up(
		db.Connection,
		"phase2-abstracting-auth-platform/migrations/",
	); err != nil {
		panic(fmt.Errorf("Failed to goose up: %v", err))
	}

	userRepo := user.UserRepository{
		Db: db.Connection,
	}

	google_provider := google.GoogleOIDCProvider(
		auth.WithClientID(cfg.GOOGLE_CLIENT_ID),
		auth.WithCallbackURI(cfg.GOOGLE_REDIRECT_URI),
		auth.WithClientSecret(cfg.GOOGLE_CLIENT_SECRET),
		auth.WithScopes([]string{"openid", "email"}),
		auth.WithDiscoverURI(cfg.GOOGLE_DISCOVER_URI),
	)

	state_manager := adapter.RedisStateStore("redis_State:6379", "", 0)
	jwt_manager := adapter.JWTTokenManager(cfg.JWT_SECRET)
	google_oidc_handler := auth.OIDCHandler{
		StateManager: state_manager,
		Provider:     google_provider,
		TokenManager: jwt_manager,
		UserManager:  &userRepo,
	}

	http.HandleFunc(
		"/",
		LoginScreenHTML,
	)

	http.HandleFunc("GET /google/auth", google_oidc_handler.Login)
	http.HandleFunc(
		`/auth/google/callback`,
		google_oidc_handler.Callback,
	)
	http.HandleFunc("/home",
		func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "/app/static_files/home.html")
		},
	)
	errChan := make(chan error, 1)
	go func() {
		if err := http.ListenAndServe(
			fmt.Sprintf(":%s", cfg.PORT),
		nil,
	); err != nil {
		errChan<-err
	}
	}()

	if err := <-errChan; err != nil {
		fmt.Println("LA")
		fmt.Println(err)
		os.Exit(0)
	}

	
}

// Directly uses `index.html` file
func LoginScreenHTML(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w,r, "/app/static_files/index.html")
}

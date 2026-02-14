package main

import (
	"net/http"

	"auth-app/internal/platform"
	"auth-app/internal/service/auth"
)

func main() {
	cfg := platform.LoadConfig()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})

	google_provider := auth.GoogleProvider{
		ClientID:     cfg.GOOGLE_CLIENT_ID,
		ClientSecret: cfg.GOOGLE_CLIENT_SECRET,
		AuthURI:      cfg.GOOGLE_AUTH_URI,
		TokenURI:     cfg.GOOGLE_TOKEN_URI,
		RedirectURI:  cfg.GOOGLE_REDIRECT_URI,
		JWKsURI:      cfg.GOOGLE_JWKS_URI,
		Scopes:       []string{"openid", "email"},
	}

	state_manager := auth.NewStateManager()
	jwt_manager := auth.JwtManager{
		Secret: cfg.JWT_SECRET,
	}
	google_oidc_handler := auth.GoogleOIDCHandler{
		StateManager:   state_manager,
		GoogleProvider: google_provider,
		JwtManager:     &jwt_manager,
	}

	http.HandleFunc("/google/auth", google_oidc_handler.Login)
	http.HandleFunc(
		`/auth/google/callback`,
		google_oidc_handler.Callback,
	)

	http.ListenAndServe(
		":8000",
		nil,
	)
}

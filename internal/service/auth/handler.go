package auth

import (
	"net/http"
	"time"
)

type GoogleOIDCHandler struct {
	GoogleProvider
	*StateManager
	*JwtManager
}

func (g GoogleOIDCHandler) Login(w http.ResponseWriter, r *http.Request) {
	state, err := g.GenerateState()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	nonce := ""
	authUrl := g.CreateAuthUrl(
		state,
		nonce,
	)

	http.Redirect(
		w,
		r,
		authUrl,
		http.StatusFound,
	)
}

func (g GoogleOIDCHandler) Callback(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	// nonce := r.URL.Query().Get("nonce")
	code := r.URL.Query().Get("code")

	if err := g.ValidateState(state); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
	}

	g_token, err := g.ExchangeCode(code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
	}

	g_claims, err := g.VerifyIdToken(g_token.IdToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	token_str, err := g.GenerateToken(
		g_claims.Subject,
		g_claims.Email,
		g_claims.Email,
		time.Duration(time.Minute*5),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	refresh_token, err := g.GenerateToken(
		g_claims.Subject,
		g_claims.Email,
		g_claims.Email,
		time.Duration(time.Minute*5),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	http.SetCookie(
		w,
		&http.Cookie{
			Name:     "access_token",
			Value:    token_str,
			SameSite: http.SameSiteStrictMode,
			Secure:   false, // development environment helloo 🙋🏽‍♂️
			HttpOnly: true,
		},
	)

	http.SetCookie(
		w,
		&http.Cookie{
			Name:     "refresh_token",
			Value:    refresh_token,
			SameSite: http.SameSiteStrictMode,
			Secure:   false, // development environment helloo 🙋🏽‍♂️
			HttpOnly: true,
		},
	)

	http.Redirect(w, r, "/", http.StatusFound)
}

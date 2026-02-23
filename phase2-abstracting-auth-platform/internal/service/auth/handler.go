package auth

import (
	"net/http"
	"time"
)

type OIDCHandler struct {
	Provider
	StateManager
	TokenManager
	UserManager
}

func (h OIDCHandler) Login(w http.ResponseWriter, r *http.Request) {
	state := h.Generate(r.Context(), h.Provider.GetName())

	authUrl, err := h.Provider.GetAuthUrl(
		state,
		"",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(
		w,
		r,
		authUrl,
		http.StatusFound,
	)
}

func (h OIDCHandler) Callback(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	code := r.URL.Query().Get("code")

	if err := h.Validate(
		r.Context(),
		state,
		h.Provider.GetName(),
	); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	// Exchange Code
	provider_token, err := h.ExchangeCode(code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	// VerifyIDToken
	provider_claims, err := h.Provider.VerifyIdToken(
		provider_token.GetIdToken(),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	provider_uid, err := provider_claims.GetSubject()
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	auth_user, err := h.UserManager.FindOrCreateUser(
		r.Context(),
		h.Provider.GetName(),
		provider_uid,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	access_tkn, err := h.GenerateToken(
		auth_user.GetId(),
		auth_user.GetEmail(),
		auth_user.GetRole(),
		"",
		"",
		time.Duration(time.Minute*5),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(
		w,
		&http.Cookie{
			Name:     "access_token",
			Value:    access_tkn,
			SameSite: http.SameSiteStrictMode,
			Secure:   false,
			HttpOnly: true,
		},
	)

	http.Redirect(w, r, "/home", http.StatusFound)
}

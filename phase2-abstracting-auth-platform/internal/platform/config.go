package platform

import "os"

type appConfig struct {
	PORT                 string `json:"port"`
	GOOGLE_CLIENT_ID     string `json:"google_client_id"`
	GOOGLE_CLIENT_SECRET string `json:"google_client_secret"`
	GOOGLE_PROJECT_ID    string `json:"google_project_id"`
	GOOGLE_X509_CERT_URL string `json:"google_x509_cert_url"`
	GOOGLE_AUTH_URI      string `json:"google_auth_uri"`
	GOOGLE_TOKEN_URI     string `json:"google_token_uri"`
	GOOGLE_REDIRECT_URI  string `json:"google_redirect_uri"`
	JWT_SECRET           string `json:"jwt_secret"`
	GOOGLE_JWKS_URI      string `json:"google_jwks_uri"`
	DB_URL               string `json:"db_url"`
	GOOGLE_DISCOVER_URI  string `json:"google_discover_uri"`
	GITHUB_CLIENT_ID     string `json:"github_client_id"`
	GITHUB_CLIENT_SECRET string `json:"github_client_secret"`
}

func LoadConfig() *appConfig {
	return &appConfig{
		PORT:                 getEnvString("PORT", "8080"),
		GOOGLE_CLIENT_ID:     getEnvString("GOOGLE_CLIENT_ID", ""),
		GOOGLE_CLIENT_SECRET: getEnvString("GOOGLE_CLIENT_SECRET", ""),
		GOOGLE_PROJECT_ID:    getEnvString("GOOGLE_PROJECT_ID", ""),
		GOOGLE_X509_CERT_URL: getEnvString("GOOGLE_X509_CERT_URL", ""),
		GOOGLE_AUTH_URI:      getEnvString("GOOGLE_AUTH_URI", ""),
		GOOGLE_TOKEN_URI:     getEnvString("GOOGLE_TOKEN_URI", ""),
		GOOGLE_REDIRECT_URI:  getEnvString("GOOGLE_REDIRECT_URI", ""),
		JWT_SECRET:           getEnvString("JWT_SECRET", ""),
		GOOGLE_JWKS_URI:      getEnvString("GOOGLE_JWKS_URI", ""),
		DB_URL:               getEnvString("DB_URL", ""),
		GOOGLE_DISCOVER_URI:  getEnvString("GOOGLE_DISCOVER_URI", "https://accounts.google.com/.well-known/openid-configuration"),
		GITHUB_CLIENT_ID:     getEnvString("GITHUB_CLIENT_ID", ""),
		GITHUB_CLIENT_SECRET: getEnvString("GITHUB_CLIENT_SECRET", ""),
	}
}

func getEnvString(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}

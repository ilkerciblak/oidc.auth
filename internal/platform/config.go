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
}

func LoadConfig() *appConfig {
	return &appConfig{
		PORT: getEnvString("PORT", "8080"),
		GOOGLE_CLIENT_ID: getEnvString("GOOGLE_CLIENT_ID", ""),
		GOOGLE_CLIENT_SECRET : getEnvString("GOOGLE_CLIENT_SECRET",""),
		GOOGLE_PROJECT_ID    : getEnvString("GOOGLE_PROJECT_ID",""),
		GOOGLE_X509_CERT_URL : getEnvString("GOOGLE_X509_CERT_URL",""),
		GOOGLE_AUTH_URI      : getEnvString("GOOGLE_AUTH_URI",""),
		GOOGLE_TOKEN_URI     : getEnvString("GOOGLE_TOKEN_URI",""),
		GOOGLE_REDIRECT_URI  : getEnvString("GOOGLE_REDIRECT_URI",""),
	}
}

func getEnvString(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}

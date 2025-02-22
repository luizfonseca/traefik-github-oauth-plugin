package traefik_github_oauth_server

import (
	"os"
	"strings"

	"github.com/spf13/cast"
)

type Config struct {
	ApiBaseURL              string
	ApiSecretKey            string
	ServerAddress           string
	DebugMode               bool
	LogLevel                string
	GitHubOAuthClientID     string
	GitHubOAuthClientSecret string
	Addr                    string
	GithubOauthScopes       []string
}

func envWithDefault(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func githubOauthScopeConfigs() []string {
	scopesFromEnv := os.Getenv("GITHUB_OAUTH_SCOPES")
	if scopesFromEnv != "" {
		return strings.Split(scopesFromEnv, ",")
	}

	return []string{}
}

func NewConfigFromEnv() *Config {
	return &Config{
		ApiBaseURL:              os.Getenv("API_BASE_URL"),
		ApiSecretKey:            os.Getenv("API_SECRET_KEY"),
		ServerAddress:           os.Getenv("SERVER_ADDRESS"),
		DebugMode:               cast.ToBool(os.Getenv("DEBUG_MODE")),
		LogLevel:                envWithDefault("LOG_LEVEL", "INFO"),
		GitHubOAuthClientID:     os.Getenv("GITHUB_OAUTH_CLIENT_ID"),
		GitHubOAuthClientSecret: os.Getenv("GITHUB_OAUTH_CLIENT_SECRET"),
		GithubOauthScopes:       githubOauthScopeConfigs(),
	}
}

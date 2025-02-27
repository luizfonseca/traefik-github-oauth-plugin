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

func envFromFile(key string) string {
	fileEnvKey := key + "_FILE"

	if value := os.Getenv(fileEnvKey); value != "" {
		content, err := os.ReadFile(value)
		if err == nil {
			return strings.TrimSpace(string(content))
		}
	}
	return ""
}

func envString(key string) string {
	if value := envFromFile(key); value != "" {
		return value
	}
	return os.Getenv(key)
}

func envWithDefault(key string, defaultValue string) string {
	value := envString(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func githubOauthScopeConfigs() []string {
	scopesFromEnv := envString("GITHUB_OAUTH_SCOPES")
	if scopesFromEnv != "" {
		return strings.Split(scopesFromEnv, ",")
	}

	return []string{}
}

func NewConfigFromEnv() *Config {
	return &Config{
		ApiBaseURL:              envString("API_BASE_URL"),
		ApiSecretKey:            envString("API_SECRET_KEY"),
		ServerAddress:           envString("SERVER_ADDRESS"),
		DebugMode:               cast.ToBool(envString("DEBUG_MODE")),
		LogLevel:                envWithDefault("LOG_LEVEL", "INFO"),
		GitHubOAuthClientID:     envString("GITHUB_OAUTH_CLIENT_ID"),
		GitHubOAuthClientSecret: envString("GITHUB_OAUTH_CLIENT_SECRET"),
		GithubOauthScopes:       githubOauthScopeConfigs(),
	}
}

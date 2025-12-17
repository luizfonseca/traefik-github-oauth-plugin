package traefik_github_oauth_plugin

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/luizfonseca/traefik-github-oauth-plugin/internal/app/traefik-github-oauth-server/model"
	"github.com/luizfonseca/traefik-github-oauth-plugin/internal/pkg/constant"
	"github.com/luizfonseca/traefik-github-oauth-plugin/internal/pkg/jwt"
	"github.com/scylladb/go-set/strset"
)

const (
	DefaultConfigAuthPath = "/_auth"
	OneDayInHours         = 24
)

// Config the middleware configuration.
type Config struct {
	ApiBaseUrl           string          `json:"api_base_url,omitempty"`
	ApiSecretKey         string          `json:"api_secret_key,omitempty"`
	AuthPath             string          `json:"auth_path,omitempty"`
	JwtSecretKey         string          `json:"jwt_secret_key,omitempty"`
	JwtExpirationInHours int64           `json:"jwt_expiration_in_hours,omitempty"`
	LogLevel             string          `json:"log_level,omitempty"`
	Whitelist            ConfigWhitelist `json:"whitelist,omitempty"`
}

// ConfigWhitelist the middleware configuration whitelist.
type ConfigWhitelist struct {
	TwoFactorAuthRequired string `json:"two_factor_auth_required,omitempty"`

	// Ids the GitHub user id list.
	Ids []string `json:"ids,omitempty"`
	// Logins the GitHub user login list.
	Logins []string `json:"logins,omitempty"`

	// Team IDs that the user must be a member of
	Teams []string `json:"teams,omitempty"`
}

// CreateConfig creates the default middleware configuration. Required by Traefik.
func CreateConfig() *Config {
	return &Config{
		ApiBaseUrl:           "",
		ApiSecretKey:         "",
		AuthPath:             DefaultConfigAuthPath,
		JwtSecretKey:         getRandomString32(),
		JwtExpirationInHours: OneDayInHours,
		Whitelist: ConfigWhitelist{
			Ids:                   []string{},
			Logins:                []string{},
			Teams:                 []string{},
			TwoFactorAuthRequired: "false",
		},
	}
}

// TraefikGithubOauthMiddleware the middleware.
type TraefikGithubOauthMiddleware struct {
	ctx  context.Context
	next http.Handler
	name string

	apiBaseUrl           string
	apiSecretKey         string
	authPath             string
	jwtSecretKey         string
	jwtExpirationInHours int64
	whitelistIdSet       *strset.Set
	whitelistLoginSet    *strset.Set
	whitelistTeamSet     *strset.Set
	whitelistRequires2FA bool

	logger *log.Logger
}

var _ http.Handler = (*TraefikGithubOauthMiddleware)(nil)

// New creates a new TraefikGithubOauthMiddleware. Required by Traefik.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	logger := log.New(os.Stdout, "service=traefik-github-oauth-middleware level=debug msg=", 0)
	// endregion Setup logger

	authPath := config.AuthPath
	if !strings.HasPrefix(authPath, "/") {
		authPath = "/" + authPath
	}

	baseUrl := strings.TrimSuffix(config.ApiBaseUrl, "/")

	return &TraefikGithubOauthMiddleware{
		ctx:  ctx,
		next: next,
		name: name,

		apiBaseUrl:           baseUrl,
		apiSecretKey:         config.ApiSecretKey,
		authPath:             authPath,
		jwtSecretKey:         config.JwtSecretKey,
		jwtExpirationInHours: config.JwtExpirationInHours,
		whitelistIdSet:       strset.New(config.Whitelist.Ids...),
		whitelistLoginSet:    strset.New(config.Whitelist.Logins...),
		whitelistTeamSet:     strset.New(config.Whitelist.Teams...),
		whitelistRequires2FA: config.Whitelist.TwoFactorAuthRequired == "true",

		logger: logger,
	}, nil
}

// ServeHTTP implements http.Handler.
func (tg *TraefikGithubOauthMiddleware) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// If the request matches the injected `/_auth` path, handle it as an authentication request.
	if req.URL.Path == tg.authPath {
		tg.handleAuthRequest(rw, req)
		return
	}

	// Otherwise, handle it as a request that has already been handled through oauth
	tg.handleRequest(rw, req)
}

// handleRequest
func (middleware *TraefikGithubOauthMiddleware) handleRequest(rw http.ResponseWriter, req *http.Request) {
	user, err := middleware.getGitHubUserFromCookie(req)
	// If cookie is missing, re-trigger oauth flow
	if err != nil {
		if req.Method == http.MethodGet {
			middleware.redirectToOAuthPage(rw, req)
			return
		}
		middleware.logger.Printf("Failed to get user from cookie: %s", err.Error())
		http.Error(rw, "", http.StatusUnauthorized)
		return
	}

	// Early check for 2FA -- if user is not whitelisted and 2FA is required, return 401
	// if middleware.whitelistRequires2FA && !user.TwoFactorEnabled {
	// 	setNoCacheHeaders(rw)
	// 	http.Error(rw, "", http.StatusUnauthorized)
	// 	return
	// }

	// If cookie is present, check if user is whitelisted
	// If nothing can be found, returns 404 as we don't want to leak information
	// But we log the error internally
	// We are also checking for the user's teams IDs
	if !middleware.whitelistIdSet.Has(user.Id) &&
		!middleware.whitelistLoginSet.Has(user.Login) && !middleware.whitelistTeamSet.HasAny(user.Teams...) {
		setNoCacheHeaders(rw)
		http.Error(rw, "", http.StatusUnauthorized)
		return
	}

	middleware.next.ServeHTTP(rw, req)
}

// handleAuthRequest
func (p TraefikGithubOauthMiddleware) handleAuthRequest(rw http.ResponseWriter, req *http.Request) {
	setNoCacheHeaders(rw)

	rid := req.URL.Query().Get(constant.QUERY_KEY_REQUEST_ID)
	result, err := p.getAuthResult(rid)
	if err != nil {
		p.logger.Printf("Failed to get auth: %s", err.Error())
		http.Error(rw, "", http.StatusInternalServerError)
		return
	}

	exp := time.Now().Add(time.Duration(p.jwtExpirationInHours) * time.Hour)

	// Generate JWTs
	tokenString, err := jwt.GenerateJwtTokenString(
		result.GitHubUserID,
		result.GitHubUserLogin,
		result.GithubTeamIDs,
		p.jwtSecretKey,
		exp,
	)
	if err != nil {
		p.logger.Printf("Failed to generate JWT: %s", err.Error())
		http.Error(rw, "", http.StatusInternalServerError)
		return
	}
	// Determine if the request is secure (HTTPS)
	secure := req.TLS != nil

	http.SetCookie(rw, &http.Cookie{
		Name:     constant.COOKIE_NAME_JWT,
		Value:    tokenString,
		Path:     "/",
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
		Expires:  exp,
	})
	http.Redirect(rw, req, result.RedirectURI, http.StatusFound)
}

// Redirects to /oauth/page-url to start oauth flow
func (p TraefikGithubOauthMiddleware) redirectToOAuthPage(rw http.ResponseWriter, req *http.Request) {
	setNoCacheHeaders(rw)

	oAuthPageURL, err := p.generateOAuthPageURL(getRawRequestUrl(req), p.getAuthURL(req))
	if err != nil {
		p.logger.Printf("Failed to generate oauth page url: %s", err.Error())
		http.Error(rw, "", http.StatusInternalServerError)
		return
	}
	http.Redirect(rw, req, oAuthPageURL, http.StatusFound)
}

func (tg TraefikGithubOauthMiddleware) generateOAuthPageURL(redirectURI, authURL string) (string, error) {
	reqBody := model.RequestGenerateOAuthPageURL{
		RedirectURI: redirectURI,
		AuthURL:     authURL,
	}

	httpClient := http.Client{Timeout: 5 * time.Second}
	reqBodyJson, err := json.Marshal(&reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, tg.getOauthPageUrl(), bytes.NewBuffer(reqBodyJson))
	if err != nil {
		return "", err
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	// If API secret key is present, use it as bearer token for internal calls
	if len(tg.apiSecretKey) > 0 {
		req.Header.Add("authorization", fmt.Sprintf("%s %s", "Bearer", tg.apiSecretKey))
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody := model.ResponseGenerateOAuthPageURL{}

	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return "", err
	}

	return respBody.OAuthPageURL, nil
}

func (tg TraefikGithubOauthMiddleware) getAuthResult(rid string) (*model.ResponseGetAuthResult, error) {
	httpClient := http.Client{Timeout: 5 * time.Second}

	req, err := http.NewRequest(http.MethodGet, tg.getOauthResultUrl(), nil)
	if err != nil {
		return nil, err
	}

	// Appends ?rid= to /oauth/result
	qr := req.URL.Query()
	qr.Add(constant.QUERY_KEY_REQUEST_ID, rid)
	req.URL.RawQuery = qr.Encode()

	if len(tg.apiSecretKey) > 0 {
		req.Header.Add("authorization", fmt.Sprintf("%s %s", "Bearer", tg.apiSecretKey))
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody := model.ResponseGetAuthResult{}
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		tg.logger.Printf("Failed to decode response from oauth server: %s", err.Error())
		return nil, err
	}

	return &respBody, nil
}

func (p *TraefikGithubOauthMiddleware) getGitHubUserFromCookie(req *http.Request) (*jwt.PayloadUser, error) {
	jwtCookie, err := req.Cookie(constant.COOKIE_NAME_JWT)
	if err != nil {
		return nil, err
	}
	return jwt.ParseTokenString(jwtCookie.Value, p.jwtSecretKey)
}

// Returns base_url + '/oauth/page-url'
func (p TraefikGithubOauthMiddleware) getOauthPageUrl() string {
	return fmt.Sprintf("%s/%s/%s", p.apiBaseUrl, constant.ROUTER_GROUP_PATH_OAUTH, constant.ROUTER_PATH_OAUTH_PAGE_URL)
}

// Returns base_url + '/oauth/result'
func (p TraefikGithubOauthMiddleware) getOauthResultUrl() string {
	return fmt.Sprintf("%s/%s/%s", p.apiBaseUrl, constant.ROUTER_GROUP_PATH_OAUTH, constant.ROUTER_PATH_OAUTH_RESULT)
}

func (tg TraefikGithubOauthMiddleware) getAuthURL(originalReq *http.Request) string {
	scheme := "http"
	if originalReq.TLS != nil {
		scheme = "https"
	}

	gen := url.URL{
		Scheme: scheme,
		Host:   originalReq.Host,
		Path:   tg.authPath,
	}

	return gen.String()
}

func setNoCacheHeaders(rw http.ResponseWriter) {
	rw.Header().Set(constant.HTTP_HEADER_CACHE_CONTROL, "no-cache, no-store, must-revalidate, private")
	rw.Header().Set(constant.HTTP_HEADER_PRAGMA, "no-cache")
	rw.Header().Set(constant.HTTP_HEADER_EXPIRES, "0")
}

func getRawRequestUrl(originalReq *http.Request) string {
	url := url.URL{}

	scheme := "http"
	if originalReq.TLS != nil {
		scheme = "https"
	}

	url.Scheme = scheme
	url.Host = originalReq.Host
	url.Path = originalReq.URL.Path
	url.RawQuery = originalReq.URL.RawQuery

	return url.String()
}

func getRandomString32() string {
	randBytes := make([]byte, 16)
	_, _ = rand.Read(randBytes)
	return hex.EncodeToString(randBytes)
}

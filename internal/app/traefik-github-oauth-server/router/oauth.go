package router

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/ggicci/httpin"
	"github.com/go-chi/render"
	"github.com/google/go-github/v49/github"
	server "github.com/luizfonseca/traefik-github-oauth-plugin/internal/app/traefik-github-oauth-server"
	"github.com/luizfonseca/traefik-github-oauth-plugin/internal/app/traefik-github-oauth-server/model"
	"github.com/luizfonseca/traefik-github-oauth-plugin/internal/pkg/constant"
	"github.com/spf13/cast"
	"golang.org/x/oauth2"
)

var (
	ErrInvalidApiBaseURL = fmt.Errorf("invalid api base url")
	ErrInvalidRID        = fmt.Errorf("invalid rid")
	ErrInvalidAuthURL    = fmt.Errorf("invalid auth url")
)

// GET /oauth/page-url
func OauthPageUrlHandler(app *server.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var reqBody *model.RequestGenerateOAuthPageURL

		err := json.NewDecoder(r.Body).Decode(&reqBody)
		if err != nil || reqBody == nil || reqBody.AuthURL == "" || reqBody.RedirectURI == "" {
			app.Logger.Error().Msgf("Missing required input params")
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, model.ResponseError{
				Message: "BadRequest",
			})
			return
		}

		rid := app.AuthRequestManager.Insert(&model.AuthRequest{
			RedirectURI: reqBody.RedirectURI,
			AuthURL:     reqBody.AuthURL,
		})

		redirectURI, err := buildRedirectURI(app.Config.ApiBaseURL, rid)
		if err != nil {
			app.Logger.Error().
				Caller().
				Stack().
				Err(err).
				Str("rid", rid).
				Msg("failed to build redirect uri")

			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, model.ResponseError{
				Message: "InternalServerError",
			})
			return
		}

		oAuthPageURL := app.GitHubOAuthConfig.AuthCodeURL(
			"",
			oauth2.SetAuthURLParam(constant.QUERY_KEY_REDIRECT_URI, redirectURI),
		)

		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, model.ResponseGenerateOAuthPageURL{
			OAuthPageURL: oAuthPageURL,
		})
	}
}

// Get /oauth/redirect
func OauthRedirectHandler(app *server.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setNoCacheHeaders(w)

		query := r.Context().Value(httpin.Input).(*model.RequestRedirect)
		if query == nil {
			app.Logger.Debug().Msg("invalid request missing RID")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("BadRequest"))
			return
		}
		authRequest, found := app.AuthRequestManager.Get(query.RID)

		if !found {
			app.Logger.Debug().Str("rid", query.RID).Msg("invalid rid")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(ErrInvalidRID.Error()))
			return
		}

		user, err := oAuthCodeToUser(r.Context(), app.GitHubOAuthConfig, query.Code)
		if err != nil {
			app.Logger.Error().
				Caller().
				Stack().
				Str("rid", query.RID).
				Str("redirect_uri", authRequest.RedirectURI).
				Str("auth_url", authRequest.AuthURL).
				Err(err).
				Msg("failed to get GitHub user")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("InternalServerError"))
			return
		}

		authRequest.GitHubUserID = cast.ToString(user.GetID())
		authRequest.GitHubUserLogin = user.GetLogin()

		authURL, err := url.Parse(authRequest.AuthURL)
		if err != nil {
			app.Logger.Warn().
				Err(err).
				Str("rid", query.RID).
				Str("auth_url", authRequest.AuthURL).
				Msg("invalid auth url")

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("InternalServerError"))
			return
		}
		authURLQuery := authURL.Query()
		authURLQuery.Set(constant.QUERY_KEY_REQUEST_ID, query.RID)
		authURL.RawQuery = authURLQuery.Encode()

		http.Redirect(w, r, authURL.String(), http.StatusFound)
	}
}

// Get /oauth/result
func OauthAuthResultHandler(app *server.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setNoCacheHeaders(w)

		query := r.Context().Value(httpin.Input).(*model.RequestGetAuthResult)
		if query == nil {
			app.Logger.Debug().Msg("invalid request missing RID")
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, model.ResponseError{
				Message: "BadRequest",
			})
			return
		}

		authRequest, found := app.AuthRequestManager.Pop(query.RID)
		if !found {
			app.Logger.Debug().Str("rid", query.RID).Msg("invalid rid")
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, model.ResponseError{
				Message: ErrInvalidRID.Error(),
			})
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(
			w,
			r,
			model.ResponseGetAuthResult{
				RedirectURI:     authRequest.RedirectURI,
				GitHubUserID:    authRequest.GitHubUserID,
				GitHubUserLogin: authRequest.GitHubUserLogin,
			},
		)
	}
}

func oAuthCodeToUser(ctx context.Context, oAuthConfig *oauth2.Config, code string) (*github.User, error) {
	ctxExchange, cancelExchange := context.WithCancel(ctx)
	defer cancelExchange()
	token, err := oAuthConfig.Exchange(ctxExchange, code)
	if err != nil {
		return nil, err
	}
	ctxClient, cancelClient := context.WithCancel(ctx)
	defer cancelClient()

	gitHubApiHttpClient := oAuthConfig.Client(ctxClient, token)
	gitHubApiClient := github.NewClient(gitHubApiHttpClient)

	ctxGetUser, cancelGetUser := context.WithCancel(ctx)
	defer cancelGetUser()

	user, _, err := gitHubApiClient.Users.Get(ctxGetUser, "")
	if err != nil {
		return nil, err
	}
	return user, nil
}

func buildRedirectURI(apiBaseUrl, rid string) (string, error) {
	redirectURI, err := url.Parse(apiBaseUrl)
	if err != nil {
		return "", ErrInvalidApiBaseURL
	}
	redirectURI = redirectURI.JoinPath(constant.ROUTER_GROUP_PATH_OAUTH, constant.ROUTER_PATH_OAUTH_REDIRECT)
	redirectURLQuery := redirectURI.Query()
	redirectURLQuery.Set(constant.QUERY_KEY_REQUEST_ID, rid)
	redirectURI.RawQuery = redirectURLQuery.Encode()
	return redirectURI.String(), nil
}

func setNoCacheHeaders(w http.ResponseWriter) {
	w.Header().Add("cache-control", "no-cache, no-store, must-revalidate, private")
	w.Header().Add("pragma", "no-cache")
	w.Header().Add("expires", "0")
}

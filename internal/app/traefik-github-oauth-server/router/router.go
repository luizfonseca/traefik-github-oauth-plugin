package router

import (
	"net/http"

	server "github.com/MuXiu1997/traefik-github-oauth-plugin/internal/app/traefik-github-oauth-server"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(app *server.App) {
	apiSecretKeyMiddleware := server.NewApiSecretKeyMiddleware(app.Config.ApiSecretKey)

	app.Router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	app.Router.Route("/oauth", func(r chi.Router) {
		r.Get("/redirect", OauthRedirectHandler(app))
		r.With(apiSecretKeyMiddleware).Post("/page-url", OauthPageUrlHandler(app))
		r.With(apiSecretKeyMiddleware).Get("/result", OauthAuthResultHandler(app))
	})
}

package router

import (
	"net/http"

	"github.com/ggicci/httpin"
	"github.com/go-chi/chi/v5"
	server "github.com/luizfonseca/traefik-github-oauth-plugin/internal/app/traefik-github-oauth-server"
	"github.com/luizfonseca/traefik-github-oauth-plugin/internal/app/traefik-github-oauth-server/model"
)

func RegisterRoutes(app *server.App) {
	apiSecretKeyMiddleware := server.NewApiSecretKeyMiddleware(app.Config.ApiSecretKey)

	app.Router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	app.Router.Route("/oauth", func(r chi.Router) {
		r.With(httpin.NewInput(model.RequestRedirect{})).Get("/redirect", OauthRedirectHandler(app))
		r.With(apiSecretKeyMiddleware).Post("/page-url", OauthPageUrlHandler(app))
		r.With(apiSecretKeyMiddleware, httpin.NewInput(model.RequestGetAuthResult{})).Get("/result", OauthAuthResultHandler(app))
	})
}

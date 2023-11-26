package traefik_github_oauth_server

import (
	"fmt"
	"net/http"

	"github.com/MuXiu1997/traefik-github-oauth-plugin/internal/app/traefik-github-oauth-server/model"
	"github.com/MuXiu1997/traefik-github-oauth-plugin/internal/pkg/constant"

	"github.com/go-chi/render"
)

// NewApiSecretKeyMiddleware returns a middleware that checks the api secret key.
func NewApiSecretKeyMiddleware(apiSecretKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// If api secret key is empty, skip the check.
			if len(apiSecretKey) == 0 {
				next.ServeHTTP(w, r)
				return
			}

			reqSecretKey := r.Header.Get("Authorization")
			if reqSecretKey != fmt.Sprintf("%s %s", constant.AUTHORIZATION_PREFIX_TOKEN, apiSecretKey) {
				render.JSON(w, r, model.ResponseError{
					Message: "Unauthorized",
				})
				next.ServeHTTP(w, r)
			}
		})
	}

}

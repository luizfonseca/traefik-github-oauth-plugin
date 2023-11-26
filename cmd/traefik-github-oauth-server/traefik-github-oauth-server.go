package main

import (
	. "github.com/luizfonseca/traefik-github-oauth-plugin/internal/app/traefik-github-oauth-server"
	"github.com/luizfonseca/traefik-github-oauth-plugin/internal/app/traefik-github-oauth-server/router"
)

func main() {
	app := NewDefaultApp()
	router.RegisterRoutes(app)
	app.Run()
}

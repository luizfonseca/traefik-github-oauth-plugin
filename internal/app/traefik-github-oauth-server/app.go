package traefik_github_oauth_server

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/patrickmn/go-cache"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"golang.org/x/oauth2"

	oauth2github "golang.org/x/oauth2/github"
)

// App the Traefik GitHub OAuth server application.
type App struct {
	Config             *Config
	Server             *http.Server
	Router             *chi.Mux
	GitHubOAuthConfig  *oauth2.Config
	AuthRequestManager *AuthRequestManager
	Logger             *zerolog.Logger
}

func NewApp(
	config *Config,
	server *http.Server,
	router *chi.Mux,
	authRequestManager *AuthRequestManager,
	logger *zerolog.Logger,
) *App {
	if config.DebugMode {
		config.LogLevel = "DEBUG"
	}

	switch config.LogLevel {
	case "DEBUG", "debug":
		logger.Level(zerolog.DebugLevel)
	case "INFO", "info":
		logger.Level(zerolog.InfoLevel)
	case "WARNING", "warning", "WARN", "warn":
		logger.Level(zerolog.WarnLevel)
	case "ERROR", "error":
		logger.Level(zerolog.ErrorLevel)
	default:
		logger.Level(zerolog.InfoLevel)
	}

	server.Addr = config.ServerAddress
	server.Handler = router

	app := &App{
		Config: config,
		Server: server,
		Router: router,
		GitHubOAuthConfig: &oauth2.Config{
			ClientID:     config.GitHubOAuthClientID,
			ClientSecret: config.GitHubOAuthClientSecret,
			Endpoint:     oauth2github.Endpoint,
			Scopes:       config.GithubOauthScopes,
		},
		AuthRequestManager: authRequestManager,
		Logger:             logger,
	}

	return app
}

func NewDefaultApp() *App {
	logger := zerolog.New(os.Stdout).With().Str("service", "traefik-github-oauth").Timestamp().Logger()

	router := chi.NewRouter()

	// Add middleware to provide more access information through logs
	router.Use(hlog.NewHandler(logger))
	router.Use(hlog.AccessHandler(
		func(r *http.Request, status, size int, duration time.Duration) {
			hlog.FromRequest(r).Info().
				Str("method", r.Method).
				Str("host", r.Host).
				Stringer("url", r.URL).
				Int("status", status).
				Int("size", size).
				Dur("duration", duration).
				Msg("")
		}),
		hlog.RefererHandler("referer"),
		hlog.UserAgentHandler("userAgent"),
		hlog.RequestIDHandler("requestId", "Request-Id"),
		hlog.RemoteAddrHandler("ip"),
		hlog.CustomHeaderHandler("xForwardedFor", "X-Forwarded-For"),
		hlog.CustomHeaderHandler("xRealIp", "X-Real-Ip"),
	)
	// Recoverer middleware recovers from panics and writes a 500 if there was one.
	router.Use(middleware.Recoverer)

	config := NewConfigFromEnv()
	return NewApp(
		config,
		&http.Server{
			ReadHeaderTimeout: 5 * time.Second,
		},
		router,
		NewAuthRequestManager(cache.New(5*time.Minute, 11*time.Minute)),
		&logger,
	)
}

func (app *App) Run() {
	app.Logger.Info().Msgf("Server listening on %s", app.Server.Addr)
	go func() {

		if err := app.Server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			app.Logger.Fatal().Err(err).Msgf("Failed to listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	app.Logger.Info().Msg("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	if err := app.Server.Shutdown(ctx); err != nil {
		app.Logger.Fatal().Err(err).Msgf("Error while shutting down server: %s\n", err)
	}
	defer cancel()
	app.Logger.Info().Msg("Server exiting")
}

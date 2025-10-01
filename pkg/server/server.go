package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ekristen/go-telemetry"
	"github.com/rs/zerolog/log"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/ekristen/go-project-template/pkg/docs"
	"github.com/ekristen/go-project-template/pkg/registry"
	"github.com/ekristen/go-project-template/pkg/router"

	_ "github.com/ekristen/go-project-template/pkg/cookies"
	_ "github.com/ekristen/go-project-template/pkg/hashes"
)

type Options struct {
	Port int

	Telemetry *telemetry.Telemetry
}

func Run(ctx context.Context, opts *Options) error {
	logger := log.With().Str("component", "server").Logger()

	r := router.Configure()

	r.Wrap(
		chimiddleware.Recoverer,
		chimiddleware.StripSlashes,
	)

	routeOpts := &registry.RouteOptions{}

	// Register all the routes, this is a nice little trick to get a chi.Router on the web.Service
	// but still allow all the fancy magic of rest service to take place.
	r.Group(func(r chi.Router) {
		for id, h := range registry.GetRegistry() {
			logger.Debug().Str("id", id).Msg("registering route")
			router.Register(r, h, routeOpts)
		}
	})

	// Setup Scalable Web UI (nicer UI)
	r.Docs("/api/v1/docs", docs.New)

	// Below this point is where the server is started and graceful shutdown occurs.

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", opts.Port),
		Handler:           r.Router,
		ReadTimeout:       1 * time.Second,
		WriteTimeout:      1 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal().Err(err).Msg("listen error")
		}
	}()
	logger.Info().Int("port", opts.Port).Msg("starting api server")

	logger.Debug().Msg("waiting for context to be done")

	<-ctx.Done()

	logger.Info().Msg("shutting down api server")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error().Err(err).Msg("unable to shutdown the api server gracefully")
		return err
	}

	return nil
}

package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ekristen/go-telemetry/v2"
	"github.com/riandyrn/otelchi"
	"github.com/sirupsen/logrus"

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
	logger := logrus.WithField("component", "server")

	r := router.Configure()

	// Add OpenTelemetry middleware for automatic HTTP instrumentation
	var middlewares []func(http.Handler) http.Handler
	middlewares = append(middlewares, chimiddleware.Recoverer)
	middlewares = append(middlewares, chimiddleware.StripSlashes)

	// Add otelchi middleware if telemetry is enabled
	if opts.Telemetry != nil && opts.Telemetry.TracerProvider() != nil {
		middlewares = append(middlewares, otelchi.Middleware(
			opts.Telemetry.ServiceName(),
			otelchi.WithTracerProvider(opts.Telemetry.TracerProvider()),
			otelchi.WithChiRoutes(r.Router),
		))
	}

	r.Wrap(middlewares...)

	routeOpts := &registry.RouteOptions{
		Telemetry: opts.Telemetry,
	}

	// Register all the routes, this is a nice little trick to get a chi.Router on the web.Service
	// but still allow all the fancy magic of rest service to take place.
	r.Group(func(r chi.Router) {
		for id, h := range registry.GetRegistry() {
			logger.WithField("id", id).Debug("registering route")
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

	// Channel to capture server errors
	serverErr := make(chan error, 1)

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.WithError(err).Error("listen error")
			serverErr <- err
		}
	}()
	logger.WithField("port", opts.Port).Info("starting api server")

	logger.Debug("waiting for context to be done")

	// Wait for either context cancellation or server error
	select {
	case <-ctx.Done():
		// Context cancelled, proceed to graceful shutdown
	case err := <-serverErr:
		// Server failed to start or crashed
		return fmt.Errorf("server error: %w", err)
	}

	logger.Info("shutting down api server")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.WithError(err).Error("unable to shutdown the api server gracefully")
		return err
	}

	return nil
}

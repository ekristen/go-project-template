package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"

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

	Log *zap.Logger
}

func RunServer(ctx context.Context, opts *Options) error {
	response.DefaultErrorResponseContentType = "application/problem+json"
	response.DefaultSuccessResponseContentType = "application/json"

	r.Wrap(
		chimiddleware.Recoverer,
		chimiddleware.StripSlashes,
	)

	routeOpts := &registry.RouteOptions{}

	// Register all the routes, this is a nice little trick to get a chi.Router on the web.Service
	// but still allow all the fancy magic of rest service to take place.
	r.Group(func(r chi.Router) {
		for id, h := range registry.GetRegistry() {
			opts.Log.Debug("registering route", zap.String("id", id))
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
			opts.Log.Fatal("listen: %s\n", zap.Error(err))
		}
	}()
	opts.Log.With(zap.Int("port", opts.Port)).Info("starting api server")

	opts.Log.Debug("waiting for context to be done")

	<-ctx.Done()

	opts.Log.Info("shutting down api server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		opts.Log.Error("unable to shutdown the api server gracefully", zap.Error(err))
		return err
	}

	return nil
}

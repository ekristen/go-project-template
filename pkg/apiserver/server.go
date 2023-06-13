package apiserver

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type Options struct {
	Port int

	Log *logrus.Entry
}

func RunServer(ctx context.Context, opts *Options) error {
	if opts.Log == nil {
		opts.Log = logrus.WithField("component", "api-server")
	} else {
		opts.Log = opts.Log.WithField("component", "api-server")
	}

	router := mux.NewRouter().StrictSlash(true)
	router.Path("/").HandlerFunc(RootHandler)

	// Below this point is where the server is started and graceful shutdown occurs.

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", opts.Port),
		Handler:           router,
		ReadTimeout:       1 * time.Second,
		WriteTimeout:      1 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			opts.Log.Fatalf("listen: %s\n", err)
		}
	}()
	opts.Log.WithField("port", opts.Port).Info("starting api server")

	<-ctx.Done()

	opts.Log.Info("shutting down api server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		opts.Log.WithError(err).Error("unable to shutdown the api server gracefully")
		return err
	}

	return nil
}

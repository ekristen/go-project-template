package apiserver

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"

	"github.com/gorilla/mux"
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
		Addr:    fmt.Sprintf(":%d", opts.Port),
		Handler: router,
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
	defer func() {
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		opts.Log.WithError(err).Error("unable to shutdown the api server gracefully")
		return err
	}

	return nil
}

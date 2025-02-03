package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/go-chi/chi/v5"

	"github.com/swaggest/jsonschema-go"
	oapi "github.com/swaggest/openapi-go"
	"github.com/swaggest/openapi-go/openapi31"
	"github.com/swaggest/rest"
	"github.com/swaggest/rest/nethttp"
	"github.com/swaggest/rest/response"
	"github.com/swaggest/rest/response/gzip"
	"github.com/swaggest/rest/web"
	swgui "github.com/swaggest/swgui/v5emb"

	"github.com/ekristen/go-project-template/pkg/common"
	"github.com/ekristen/go-project-template/pkg/docs"
	"github.com/ekristen/go-project-template/pkg/hashes"
)

type Options struct {
	Port int

	Log *zap.Logger
}

func RunServer(ctx context.Context, opts *Options) error {
	response.DefaultErrorResponseContentType = "application/problem+json"
	response.DefaultSuccessResponseContentType = "application/json"

	r := openapi31.NewReflector()
	s := web.NewService(r)

	s.OpenAPISchema().SetTitle(common.NAME)
	s.OpenAPISchema().SetDescription(common.NAME)
	s.OpenAPISchema().SetVersion(common.VERSION)

	// An example of global schema override to disable additionalProperties for all object schemas.
	jsr := s.OpenAPIReflector().JSONSchemaReflector()
	jsr.DefaultOptions = append(jsr.DefaultOptions,
		jsonschema.InterceptSchema(func(params jsonschema.InterceptSchemaParams) (stop bool, err error) {
			// Allow unknown request headers and skip response.
			if oc, ok := oapi.OperationCtx(params.Context); !params.Processed || !ok ||
				oc.IsProcessingResponse() || oc.ProcessingIn() == oapi.InHeader {
				return false, nil
			}

			schema := params.Schema

			if schema.HasType(jsonschema.Object) && len(schema.Properties) > 0 && schema.AdditionalProperties == nil {
				schema.AdditionalProperties = (&jsonschema.SchemaOrBool{}).WithTypeBoolean(false)
			}

			return false, nil
		}),
	)
	s.OpenAPICollector.CombineErrors = "anyOf"
	s.AddHeadToGet = true

	s.Wrap(
		// Response validator setup.
		//
		// It might be a good idea to disable this middleware in production to save performance,
		// but keep it enabled in dev/test/staging environments to catch logical issues.
		response.ValidatorMiddleware(s.ResponseValidatorFactory),
		gzip.Middleware, // Response compression with support for direct gzip pass through.

		// Example middleware to setup custom error responses.
		func(handler http.Handler) http.Handler {
			var h *nethttp.Handler
			if nethttp.HandlerAs(handler, &h) {
				h.MakeErrResp = func(ctx context.Context, err error) (int, interface{}) {
					code, er := rest.Err(err)

					var ae anotherErr

					if errors.As(err, &ae) {
						return http.StatusBadRequest, ae
					}

					return code, customErr{
						Message: er.ErrorText,
						Details: er.Context,
					}
				}
			}

			return handler
		},
	)

	// Root Handler NOT included in OpenAPI Spec
	// Example of a RAW HTTP Handler
	s.Router.Get("/", RootHandler)

	s.Route("/api/v1", func(r chi.Router) {
		hashes.Register(r)
	})

	// Setup Scalable Web UI (nicer UI)
	s.Docs("/api/v1/docs", docs.New)
	// Setup Swagger UI.
	s.Docs("/api/v1/docs/swagger", swgui.New)

	// Below this point is where the server is started and graceful shutdown occurs.

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", opts.Port),
		Handler:           s.Router,
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

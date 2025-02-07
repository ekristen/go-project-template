package router

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	oapi "github.com/swaggest/openapi-go"
	"github.com/swaggest/openapi-go/openapi31"
	"github.com/swaggest/rest/nethttp"
	"github.com/swaggest/rest/response"
	"github.com/swaggest/rest/response/gzip"
	"github.com/swaggest/rest/web"

	"github.com/ekristen/go-project-template/pkg/common"
	"github.com/ekristen/go-project-template/pkg/registry"
)

func Configure() *web.Service {
	response.DefaultErrorResponseContentType = "application/problem+json"
	response.DefaultSuccessResponseContentType = "application/json"

	r := openapi31.NewReflector()

	s := web.NewService(r)

	s.OpenAPISchema().SetTitle(common.NAME)
	s.OpenAPISchema().SetDescription("project template")
	s.OpenAPISchema().SetVersion(common.VERSION)

	s.OpenAPISchema().
		SetAPIKeySecurity("header", "x-api-key", oapi.InHeader, "api key")
	s.OpenAPISchema().
		SetAPIKeySecurity("cookie", common.NAME, oapi.InCookie, "session cookie")

	// An example of global schema override to disable additionalProperties for all object schemas.
	// TODO: this causes a bug with requests if additional cookies are sent potentially, might want to disable entirely?
	// that is why it is currently disabled.
	/*
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
	*/

	s.OpenAPICollector.CombineErrors = "anyOf"
	s.AddHeadToGet = true

	s.Wrap(
		// Response validator setup.
		//
		// It might be a good idea to disable this middleware in production to save performance,
		// but keep it enabled in dev/test/staging environments to catch logical issues.
		// TODO: re-enable after bug with schema gen is fixed ...
		//response.ValidatorMiddleware(s.ResponseValidatorFactory),
		gzip.Middleware, // Response compression with support for direct gzip pass through.
	)

	return s
}

func Register(r chi.Router, h registry.WithID, opts *registry.RouteOptions) {
	if h, ok := h.(registry.WithSetOpts); ok {
		h.SetOpts(opts)
	}

	var with chi.Middlewares
	if w, ok := h.(registry.WithWithMiddleware); ok {
		with = w.WithMiddleware()
	}

	if w, ok := h.(registry.WithPermission); ok {
		with = append(with, attachPermissions([]string{w.Permission()}))
	}

	var ucOptions []func(h *nethttp.Handler)
	if w, ok := h.(registry.WithUseCaseOptions); ok {
		ucOptions = w.UseCaseOptions()
	}

	if uc, ok := h.(registry.WithUseCase); ok {
		r.With(with...).Method(h.Method(), h.Path(), nethttp.NewHandler(uc.UseCase(), ucOptions...))
	} else if handler, ok := h.(registry.WithServeHTTP); ok {
		if h.Method() == "ALL" {
			r.With(with...).Handle(h.Path(), http.HandlerFunc(handler.ServeHTTP))
			return
		}

		r.With(with...).Method(h.Method(), h.Path(), http.HandlerFunc(handler.ServeHTTP))
	} else {
		panic("unknown handler type")
	}
}

// attachPermissions attaches permissions to the request context which allows it to be checked via
// the authorization middleware.
func attachPermissions(permissions []string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("called here")
			ctx := context.WithValue(r.Context(), "requiredPermissions", permissions)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

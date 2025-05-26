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

	s.OpenAPICollector.CombineErrors = "anyOf"
	s.AddHeadToGet = true

	s.Wrap(
		// Response validator setup.
		//
		// It might be a good idea to disable this middleware in production to save performance,
		// but keep it enabled in dev/test/staging environments to catch logical issues.
		response.ValidatorMiddleware(s.ResponseValidatorFactory),
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

	switch v := h.(type) {
	case registry.WithUseCase:
		r.With(with...).Method(h.Method(), h.Path(), nethttp.NewHandler(v.UseCase(), ucOptions...))
	case registry.WithServeHTTP:
		if h.Method() == "ALL" {
			r.With(with...).Handle(h.Path(), http.HandlerFunc(v.ServeHTTP))
			return
		}

		r.With(with...).Method(h.Method(), h.Path(), http.HandlerFunc(v.ServeHTTP))
	default:
		panic(fmt.Sprintf("unknown handler type %T", v))
	}
}

// attachPermissions attaches permissions to the request context which allows it to be checked via
// the authorization middleware.
func attachPermissions(permissions []string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), "requiredPermissions", permissions) //nolint:staticcheck
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

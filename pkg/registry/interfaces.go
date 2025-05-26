package registry

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/swaggest/rest/nethttp"
	"github.com/swaggest/usecase"
)

// WithID is the base interface for all handlers.
type WithID interface {
	ID() string
	Method() string
	Path() string
}

// WithSetOpts is an interface that allows setting route options.
type WithSetOpts interface {
	SetOpts(opts *RouteOptions)
}

// WithUseCase is an interface that allows obtaining a use case.
type WithUseCase interface {
	UseCase() usecase.Interactor
}

// WithUseCaseOptions is an interface that allows obtaining use case options to modify the use case.
type WithUseCaseOptions interface {
	UseCaseOptions() []func(h *nethttp.Handler)
}

// WithPermission is an interface that allows obtaining the permission associated with the handler.
type WithPermission interface {
	Permission() string
}

// WithWithMiddleware is an interface that allows obtaining the middleware associated with the handler.
type WithWithMiddleware interface {
	WithMiddleware() chi.Middlewares
}

// WithServeHTTP is an interface that allows serving HTTP requests.
type WithServeHTTP interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

package server

import (
	"fmt"
	"net/http"

	"github.com/ekristen/go-telemetry"

	"github.com/ekristen/go-project-template/pkg/common"
	"github.com/ekristen/go-project-template/pkg/registry"
)

func init() {
	registry.Register(&RootHandler{})
}

// RootHandler is an example of how to implement a handler without using UseCase. It's a simple handler that returns the
// name and version of the application as a standard http.Handler implementation.
type RootHandler struct {
	telemetry *telemetry.Telemetry
}

func (h *RootHandler) ID() string {
	return "root"
}

func (h *RootHandler) Method() string {
	return http.MethodGet
}

func (h *RootHandler) Path() string {
	return "/"
}

func (h *RootHandler) SetOpts(opts *registry.RouteOptions) {
	h.telemetry = opts.Telemetry
}

func (h *RootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, span, logger := h.telemetry.StartSpanWithLogger(r.Context(), "server.root")
	defer span.End()

	logger.Info().Msg("serving root endpoint")

	data := fmt.Sprintf(`{"name":%q,"version":%q}`, common.AppVersion.Name, common.AppVersion.Summary)

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write([]byte(data)); err != nil {
		logger.Warn().Err(err).Msg("unable to write to response")
		span.RecordError(err)
	} else {
		logger.Info().
			Str("app_name", common.AppVersion.Name).
			Str("app_version", common.AppVersion.Summary).
			Msg("root endpoint served successfully")
	}
}

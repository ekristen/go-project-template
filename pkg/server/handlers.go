package server

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/ekristen/go-project-template/pkg/common"
	"github.com/ekristen/go-project-template/pkg/registry"
)

func init() {
	registry.Register(&RootHandler{})
}

// RootHandler is an example of how to implement a handler without using UseCase. It's a simple handler that returns the
// name and version of the application as a standard http.Handler implementation.
type RootHandler struct{}

func (h *RootHandler) ID() string {
	return "root"
}

func (h *RootHandler) Method() string {
	return http.MethodGet
}

func (h *RootHandler) Path() string {
	return "/"
}

func (h *RootHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	data := fmt.Sprintf(`{"name":%q,"version":%q}`, common.AppVersion.Name, common.AppVersion.Summary)

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write([]byte(data)); err != nil {
		zap.L().Warn("unable to write to response", zap.Error(err))
	}
}

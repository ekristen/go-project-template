package cookies

import (
	"context"
	"net/http"

	"github.com/ekristen/go-telemetry/v2"
	"github.com/sirupsen/logrus"
	"github.com/swaggest/usecase"

	"github.com/ekristen/go-project-template/pkg/registry"
)

func init() {
	registry.Register(&ClearHandler{})
}

type ClearRequest struct {
}

type ClearResponse struct {
	SessionID string `json:"-" cookie:"sid,httponly,secure,max-age=-1,samesite=lax,path=/" description:"The session ID"`
}

type ClearHandler struct {
	telemetry *telemetry.Telemetry
}

func (h *ClearHandler) ID() string {
	return "cookies-clear"
}

func (h *ClearHandler) Method() string {
	return http.MethodGet
}

func (h *ClearHandler) Path() string {
	return "/cookies/clear"
}

func (h *ClearHandler) SetOpts(opts *registry.RouteOptions) {
	h.telemetry = opts.Telemetry
}

func (h *ClearHandler) UseCase() usecase.Interactor {
	u := usecase.NewInteractor(h.interact)

	u.SetTitle("Clear Cookie")
	u.SetName("clear-cookie")
	u.SetDescription("clear a cookie")
	u.SetTags("Cookies")

	return u
}

func (h *ClearHandler) interact(ctx context.Context, _ ClearRequest, output *ClearResponse) error {
	ctx, span := h.telemetry.StartSpan(ctx, "cookies.clear")
	defer span.End()

	// Logger will automatically pick up trace context from the span
	logger := logrus.WithContext(ctx).WithField("component", "cookies.clear")

	logger.Info("clearing cookie")

	// Note: this isn't necessary, but it's just here for some content.
	// The max-age=-1 is what actually deletes the cookie.
	output.SessionID = "delete"

	logger.WithField("session_id", output.SessionID).Info("cookie cleared successfully")

	return nil
}

package cookies

import (
	"context"
	"net/http"

	"github.com/ekristen/go-telemetry"
	"github.com/gofrs/uuid/v5"
	"github.com/swaggest/usecase"

	"github.com/ekristen/go-project-template/pkg/registry"
)

func init() {
	registry.Register(&SetHandler{})
}

type SetRequest struct {
}

type SetResponse struct {
	SessionID string `json:"-" cookie:"sid,httponly,secure,max-age=86400,samesite=lax,path=/" description:"The session ID"`
}

type SetHandler struct {
	telemetry *telemetry.Telemetry
}

func (h *SetHandler) ID() string {
	return "cookies-set"
}

func (h *SetHandler) Method() string {
	return http.MethodGet
}

func (h *SetHandler) Path() string {
	return "/cookies/set"
}

func (h *SetHandler) SetOpts(opts *registry.RouteOptions) {
	h.telemetry = opts.Telemetry
}

func (h *SetHandler) UseCase() usecase.Interactor {
	u := usecase.NewInteractor(h.interact)

	u.SetTitle("Set Cookie")
	u.SetName("set-cookie")
	u.SetDescription("set a cookie")
	u.SetTags("Cookies")

	return u
}

func (h *SetHandler) interact(ctx context.Context, _ SetRequest, output *SetResponse) error {
	spanCtx, span, logger := h.telemetry.StartSpanWithLogger(ctx, "cookies.set")
	defer span.End()
	ctx = spanCtx

	logger.Info().Msg("setting cookie")

	output.SessionID = uuid.Must(uuid.NewV7()).String()

	logger.Info().Str("session_id", output.SessionID).Msg("cookie set successfully")

	return nil
}

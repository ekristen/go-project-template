package cookies

import (
	"context"
	"net/http"

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
	SessionID string `json:"-" cookie:"gpt-sid,httponly,secure,max-age=86400,samesite=lax,path=/" description:"The session ID"`
}

type SetHandler struct{}

func (h *SetHandler) ID() string {
	return "cookies-set"
}

func (h *SetHandler) Method() string {
	return http.MethodGet
}

func (h *SetHandler) Path() string {
	return "/cookies/set"
}

func (h *SetHandler) UseCase() usecase.Interactor {
	u := usecase.NewInteractor(h.interact)

	u.SetTitle("Set Cookie")
	u.SetName("set-cookie")
	u.SetDescription("set a cookie")
	u.SetTags("Cookies")

	return u
}

func (h *SetHandler) interact(_ context.Context, _ SetRequest, output *SetResponse) error {
	output.SessionID = uuid.Must(uuid.NewV7()).String()

	return nil
}

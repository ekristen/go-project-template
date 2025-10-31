package hashes

import (
	"context"
	"net/http"

	"github.com/ekristen/go-telemetry/v2"
	"github.com/rs/zerolog/log"
	"github.com/swaggest/usecase"

	"github.com/ekristen/go-project-template/pkg/api"
	"github.com/ekristen/go-project-template/pkg/registry"
)

func init() {
	registry.Register(&IdentifyHandler{})
}

type IdentifyRequest struct {
	api.Request
	Hash string `path:"hash" example:"1234567890abcdef" description:"The hash to identify"`
}

type IdentifyResponse struct {
	Hash string `json:"hash" example:"1234567890abcdef" description:"The hash"`
	Type string `json:"type" example:"md5" description:"The type of hash"`
}

type IdentifyHandler struct {
	telemetry *telemetry.Telemetry
}

func (h *IdentifyHandler) ID() string {
	return "hash-identify"
}

func (h *IdentifyHandler) Method() string {
	return http.MethodGet
}

func (h *IdentifyHandler) Path() string {
	return "/api/v1/hash/{hash}"
}

func (h *IdentifyHandler) SetOpts(opts *registry.RouteOptions) {
	h.telemetry = opts.Telemetry
}

func (h *IdentifyHandler) UseCase() usecase.Interactor {
	u := usecase.NewInteractor(h.interact)

	u.SetTitle("Identify Hash")
	u.SetName("identify-hash")
	u.SetDescription(
		`Attempt to identify a hash based on it's length`)
	u.SetTags("Hashes")

	return u
}

func (h *IdentifyHandler) interact(ctx context.Context, input IdentifyRequest, output *IdentifyResponse) error {
	ctx, span := h.telemetry.StartSpan(ctx, "hashes.identify")
	defer span.End()

	// Logger will automatically pick up trace context from the span
	logger := log.With().Str("component", "hashes.identify").Logger()

	logger.Info().
		Str("hash", input.Hash).
		Int("hash_length", len(input.Hash)).
		Msg("identifying hash type")

	output.Hash = input.Hash
	switch len(input.Hash) {
	case 32:
		output.Type = "md5"
	case 40:
		output.Type = "sha1"
	case 64:
		output.Type = "sha256"
	default:
		output.Type = "unknown"
	}

	logger.Info().
		Str("hash", input.Hash).
		Str("identified_type", output.Type).
		Msg("hash type identified successfully")

	return nil
}

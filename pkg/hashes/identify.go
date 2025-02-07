package hashes

import (
	"context"
	"net/http"

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

func (h *IdentifyHandler) UseCase() usecase.Interactor {
	u := usecase.NewInteractor(h.interact)

	u.SetTitle("Identify Hash")
	u.SetName("identify-hash")
	u.SetDescription(
		`Attempt to identify a hash based on it's length`)
	u.SetTags("Hashes")

	return u
}

func (h *IdentifyHandler) interact(_ context.Context, input IdentifyRequest, output *IdentifyResponse) error {
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

	return nil
}

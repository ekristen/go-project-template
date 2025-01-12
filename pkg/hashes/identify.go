package hashes

import (
	"context"

	"github.com/swaggest/usecase"

	"github.com/ekristen/go-project-template/pkg/api"
)

type HashRequest struct {
	api.Request
	Hash string `path:"hash" example:"1234567890abcdef" description:"The hash to identify"`
}

type HashResponseData struct {
	Hash string `json:"hash" example:"1234567890abcdef" description:"The hash"`
	Type string `json:"type" example:"md5" description:"The type of hash"`
}

func Identify() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input HashRequest, output *HashResponseData) error {
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
	})

	u.SetTitle("Identify Hash")
	u.SetName("identify-hash")
	u.SetDescription(
		`Attempt to identify a hash based on it's length`)
	u.SetTags("Hashes")

	return u
}

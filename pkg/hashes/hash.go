package hashes

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/swaggest/usecase"

	"github.com/ekristen/go-project-template/pkg/registry"
)

func init() {
	registry.Register(&FileHandler{})
}

type FileRequest struct {
	File multipart.File `formData:"file" description:"The file to hash" required:"true"`
}

type FileResponseData struct {
	Hash string `json:"hash" example:"1234567890abcdef" description:"The hash"`
}

type FileHandler struct {
	DB interface{}
}

func (h *FileHandler) ID() string {
	return "hash-file"
}

func (h *FileHandler) Method() string {
	return http.MethodPost
}

func (h *FileHandler) Path() string {
	return "/api/v1/hash"
}

func (h *FileHandler) SetOpts(opts *registry.RouteOptions) {
	h.DB = opts.DB
}

func (h *FileHandler) UseCase() usecase.Interactor {
	u := usecase.NewInteractor(h.interact)

	u.SetTitle("Hash File")
	u.SetName("hash-file")
	u.SetDescription(
		`Upload a file to get it's hash`)
	u.SetTags("Hashes")

	return u
}

func (h *FileHandler) interact(_ context.Context, input FileRequest, output *FileResponseData) error {
	defer input.File.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, input.File); err != nil {
		return err
	}

	hash := hasher.Sum(nil)
	output.Hash = hex.EncodeToString(hash)
	return nil
}

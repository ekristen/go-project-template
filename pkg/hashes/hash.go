package hashes

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/ekristen/go-telemetry/v2"
	"github.com/sirupsen/logrus"
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
	DB        interface{}
	telemetry *telemetry.Telemetry
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
	h.telemetry = opts.Telemetry
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

func (h *FileHandler) interact(ctx context.Context, input FileRequest, output *FileResponseData) error {
	ctx, span := h.telemetry.StartSpan(ctx, "hashes.file")
	defer span.End()

	defer input.File.Close()

	// Logger will automatically pick up trace context from the span
	logger := logrus.WithContext(ctx).WithField("component", "hashes.file")

	logger.Info("hashing file")

	hasher := sha256.New()
	span.AddEvent("sha256.New")

	if _, err := io.Copy(hasher, input.File); err != nil {
		return err
	}
	span.AddEvent("io.copy")

	hash := hasher.Sum(nil)
	span.AddEvent("sha256.Sum")
	output.Hash = hex.EncodeToString(hash)

	logger.WithField("hash", output.Hash).Info("file hashed successfully")

	return nil
}

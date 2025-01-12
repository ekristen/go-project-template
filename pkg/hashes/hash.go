package hashes

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"mime/multipart"

	"github.com/swaggest/usecase"
)

type FileRequest struct {
	File multipart.File `formData:"file" description:"The file to hash" required:"true"`
}

type FileResponseData struct {
	Hash string `json:"hash" example:"1234567890abcdef" description:"The hash"`
}

func HashFile() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input FileRequest, output *FileResponseData) error {
		defer input.File.Close()

		hasher := sha256.New()
		if _, err := io.Copy(hasher, input.File); err != nil {
			return err
		}

		hash := hasher.Sum(nil)
		output.Hash = hex.EncodeToString(hash)
		return nil
	})

	u.SetTitle("Hash File")
	u.SetName("hash-file")
	u.SetDescription(
		`Upload a file to get it's hash`)
	u.SetTags("Hashes")

	return u
}

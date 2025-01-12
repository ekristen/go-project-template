package server

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/ekristen/go-project-template/pkg/common"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	data := fmt.Sprintf(`{"name":%q,"version":%q}`, common.AppVersion.Name, common.AppVersion.Summary)

	w.WriteHeader(200)
	if _, err := w.Write([]byte(data)); err != nil {
		zap.L().Warn("unable to write to response", zap.Error(err))
	}
}

package apiserver

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/swade1987/go-project-template/pkg/common"
)

var logger = zap.L()

func RootHandler(w http.ResponseWriter, r *http.Request) {
	data := fmt.Sprintf(`{"name":%q,"version":%q}`, common.AppVersion.Name, common.AppVersion.Summary)

	w.WriteHeader(200)
	if _, err := w.Write([]byte(data)); err != nil {
		logger.Warn("unable to write to response", zap.Error(err))
	}
}

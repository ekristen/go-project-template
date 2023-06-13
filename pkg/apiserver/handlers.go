package apiserver

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/ekristen/go-project-template/pkg/common"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	data := fmt.Sprintf(`{"name":%q,"version":%q}`, common.AppVersion.Name, common.AppVersion.Summary)

	w.WriteHeader(200)
	if _, err := w.Write([]byte(data)); err != nil {
		logrus.WithError(err).Warn("unable to write to response")
	}
}

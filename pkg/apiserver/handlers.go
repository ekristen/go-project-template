package apiserver

import (
	"fmt"
	"net/http"

	"github.com/ekristen/go-project-template/pkg/common"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	data := fmt.Sprintf(`{"name":"%s","version":"%s"}`, common.AppVersion.Name, common.AppVersion.Summary)

	w.WriteHeader(200)
	w.Write([]byte(data))
	return
}

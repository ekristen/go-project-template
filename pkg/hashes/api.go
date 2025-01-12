package hashes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/swaggest/rest/nethttp"
)

func Register(r chi.Router) {
	r.Method(http.MethodGet, "/hash/{hash}", nethttp.NewHandler(Identify()))
	r.Method(http.MethodPost, "/hash", nethttp.NewHandler(HashFile()))
}

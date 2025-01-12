package docs

import (
	"fmt"
	"net/http"
)

func New(title string, swaggerJSONPath string, basePath string) http.Handler {
	return &ScalarGUI{
		Title:       title,
		SwaggerJSON: swaggerJSONPath,
		BasePath:    basePath,
	}
}

type ScalarGUI struct {
	Title       string
	SwaggerJSON string
	BasePath    string
}

func (s *ScalarGUI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(fmt.Sprintf(`<!doctype html>
<html>
  <head>
    <title>%s</title>
    <meta charset="utf-8" />
    <meta
      name="viewport"
      content="width=device-width, initial-scale=1" />
  </head>
  <body>
    <script
      id="api-reference"
      data-url="%s"></script>
    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
  </body>
</html>`, s.Title, s.SwaggerJSON)))
}

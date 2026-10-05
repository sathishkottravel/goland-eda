// Package docs serves the OpenAPI spec and a Swagger UI page.
package docs

import (
	_ "embed"
	"net/http"
)

//go:embed openapi.yaml
var spec []byte

//go:embed swagger.html
var page []byte

// Register mounts /openapi.yaml and /swagger on mux.
func Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /openapi.yaml", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/yaml")
		_, _ = w.Write(spec)
	})
	mux.HandleFunc("GET /swagger", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(page)
	})
}

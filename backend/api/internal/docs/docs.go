// Package docs serves the OpenAPI spec, Swagger UI and the GraphQL playground.
package docs

import (
	_ "embed"
	"net/http"
)

//go:embed openapi.yaml
var spec []byte

//go:embed swagger.html
var swaggerPage []byte

//go:embed playground.html
var playgroundPage []byte

// Register mounts /openapi.yaml, /swagger and /playground on mux.
func Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /openapi.yaml", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/yaml")
		_, _ = w.Write(spec)
	})
	mux.HandleFunc("GET /swagger", servePage(swaggerPage))
	mux.HandleFunc("GET /playground", servePage(playgroundPage))
}

func servePage(b []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(b)
	}
}

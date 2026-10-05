// Package rest exposes the service layer as JSON over HTTP.
package rest

import (
	"encoding/json"
	"net/http"

	"github.com/sathishkottravel/goland-eda/backend/api/internal/service"
)

type Handler struct {
	greeter *service.Greeter
}

func New(g *service.Greeter) *Handler { return &Handler{greeter: g} }

// Register mounts the REST routes on mux.
func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /health", h.health)
	mux.HandleFunc("GET /api/v1/hello", h.hello)
}

func (h *Handler) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) hello(w http.ResponseWriter, r *http.Request) {
	msg := h.greeter.Hello(r.URL.Query().Get("name"))
	writeJSON(w, http.StatusOK, map[string]string{"message": msg})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

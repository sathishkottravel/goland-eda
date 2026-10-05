package docs

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRoutes(t *testing.T) {
	mux := http.NewServeMux()
	Register(mux)

	tests := []struct {
		path, contentType string
	}{
		{"/openapi.yaml", "application/yaml"},
		{"/swagger", "text/html"},
		{"/playground", "text/html"},
	}
	for _, tc := range tests {
		t.Run(tc.path, func(t *testing.T) {
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, tc.path, nil))

			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d, want 200", rec.Code)
			}
			if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, tc.contentType) {
				t.Errorf("Content-Type = %q, want %q", ct, tc.contentType)
			}
			if rec.Body.Len() == 0 {
				t.Error("empty body")
			}
		})
	}
}

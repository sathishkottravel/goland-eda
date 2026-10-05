package rest

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sathishkottravel/goland-eda/backend/api/internal/service"
)

func TestRoutes(t *testing.T) {
	mux := http.NewServeMux()
	New(service.NewGreeter()).Register(mux)

	tests := []struct {
		path, key, want string
	}{
		{"/health", "status", "ok"},
		{"/api/v1/hello", "message", "Hello, world!"},
		{"/api/v1/hello?name=Kafka", "message", "Hello, Kafka!"},
	}
	for _, tc := range tests {
		t.Run(tc.path, func(t *testing.T) {
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, tc.path, nil))

			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d, want 200", rec.Code)
			}
			var body map[string]string
			if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
				t.Fatal(err)
			}
			if body[tc.key] != tc.want {
				t.Errorf("%s = %q, want %q", tc.key, body[tc.key], tc.want)
			}
		})
	}
}

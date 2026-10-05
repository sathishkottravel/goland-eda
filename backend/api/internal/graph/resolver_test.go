package graph

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/sathishkottravel/goland-eda/backend/api/internal/service"
)

func TestQuery(t *testing.T) {
	schema := NewSchema(service.NewGreeter())

	res := schema.Exec(context.Background(), `{ health hello(name: "Kafka") }`, "", nil)
	if len(res.Errors) > 0 {
		t.Fatal(res.Errors)
	}
	var got struct{ Health, Hello string }
	if err := json.Unmarshal(res.Data, &got); err != nil {
		t.Fatal(err)
	}
	if got.Health != "ok" || got.Hello != "Hello, Kafka!" {
		t.Errorf("got %+v", got)
	}
}

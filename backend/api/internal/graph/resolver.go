// Package graph exposes the service layer over GraphQL.
package graph

import (
	_ "embed"
	"net/http"

	graphql "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"

	"github.com/sathishkottravel/goland-eda/backend/api/internal/service"
)

//go:embed schema.graphql
var schemaSDL string

type Resolver struct {
	greeter *service.Greeter
}

func (r *Resolver) Health() string { return "ok" }

func (r *Resolver) Hello(args struct{ Name *string }) string {
	name := ""
	if args.Name != nil {
		name = *args.Name
	}
	return r.greeter.Hello(name)
}

// NewSchema parses the embedded schema and binds it to the resolvers.
func NewSchema(g *service.Greeter) *graphql.Schema {
	return graphql.MustParseSchema(schemaSDL, &Resolver{greeter: g})
}

// Handler serves GraphQL queries over HTTP POST.
func Handler(s *graphql.Schema) http.Handler { return &relay.Handler{Schema: s} }

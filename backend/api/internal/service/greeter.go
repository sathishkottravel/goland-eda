// Package service holds business logic shared by the REST and GraphQL layers.
package service

import "strings"

type Greeter struct{}

func NewGreeter() *Greeter { return &Greeter{} }

// Hello returns a greeting for name, defaulting to "world".
func (g *Greeter) Hello(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		name = "world"
	}
	return "Hello, " + name + "!"
}

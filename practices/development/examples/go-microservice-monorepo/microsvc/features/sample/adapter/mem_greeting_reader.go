package adapter

import (
	"context"
	"math/rand"
)

// InMemGreetingReader implements sample.GreetingReader with an in-memory set of greetings
type InMemGreetingReader struct{}

var greetings = []string{
	"Hello",
	"Bonjour",
	"Hola",
	"Howdy",
	"Greetings",
	"Howdy-do",
}

// List implements sample.GreetingReader for InMemGreetingReader
func (rdr InMemGreetingReader) List(context.Context) ([]string, error) {
	return greetings, nil
}

// RandomGreeting implements sample.GreetingReader for InMemGreetingReader
func (rdr InMemGreetingReader) RandomGreeting(context.Context) (string, error) {
	greetingIdx := rand.Intn(len(greetings))
	return greetings[greetingIdx], nil
}

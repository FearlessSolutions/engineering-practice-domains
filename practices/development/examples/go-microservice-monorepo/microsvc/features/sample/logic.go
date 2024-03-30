package sample

import (
	"context"
	"errors"
	"example.com/sample/commonlib/logger"
	"fmt"
	"go.uber.org/zap"
	"slices"
)

// The following comment instructs the "go generate" command to run mockgen to generate our test mocks
//go:generate mockgen -destination ./logic_mocks.go -package sample . GreetingReader,GreetingWriter,Core

// GreetingReader is a driven port for something that reads greetings. It is "driven" by the business logic,
// and represents information that is fetched from outside the application such as a database, message queue, or
// HTTP connection to another microservice
type GreetingReader interface {
	// RandomGreeting retrieves a random greeting prefix
	RandomGreeting(ctx context.Context) (string, error)
	// List provides a list of all greeting prefixes available
	List(ctx context.Context) ([]string, error)
}

// GreetingWriter is a driven port for something that writes to the collection of greetings. It is "driven" by the business
// logic, and represents information that is fetched from outside the application
type GreetingWriter interface {
	// AddGreeting adds a new greeting to the set of available greetings
	AddGreeting(ctx context.Context, newGreeting string) error
}

// ErrGreetingAlreadyExists represents a situation where someone tries to add a duplicate greeting to the set of greetings
var ErrGreetingAlreadyExists = errors.New("the passed greeting already exists")

// Core is a driving port describing the capabilities of the sample feature. It separates the business logic from
// the "driving" external interface, such as a REST controller, which can then be tested in isolation from the business logic
type Core interface {
	// GiveGreeting produces a greeting which it fetches from its GreetingReader
	GiveGreeting(ctx context.Context, name string, greetingReader GreetingReader) (string, error)
	// AddGreeting adds a new greeting to the set of greetings used in GiveGreeting. It returns ErrGreetingAlreadyExists if
	// trying to insert a duplicate.
	AddGreeting(ctx context.Context, newGreeting string, greetingReader GreetingReader, greetingWriter GreetingWriter) error
}

// CoreLogic implements the core business logic of the sample feature which is plugged into the input adapter of
// the hexagonal architecture (i.e. a REST controller). Business logic cores can talk to each other directly without an interface.
type CoreLogic struct{}

// GiveGreeting implements Core for CoreLogic
func (CoreLogic) GiveGreeting(ctx context.Context, name string, greetingReader GreetingReader) (string, error) {
	// Fetch a greeting
	nextGreeting, greetingErr := greetingReader.RandomGreeting(ctx)
	if greetingErr != nil {
		return "", fmt.Errorf("could not retrieve a greeting: %w", greetingErr)
	}

	logger.Log.Debug("Fetched a greeting.", zap.String("greeting", nextGreeting))
	return fmt.Sprintf("%v, %v!", nextGreeting, name), nil
}

// AddGreeting implements Core for CoreLogic
func (CoreLogic) AddGreeting(ctx context.Context, newGreeting string, greetingReader GreetingReader, greetingWriter GreetingWriter) error {
	greetings, greetingReadErr := greetingReader.List(ctx)
	if greetingReadErr != nil {
		return fmt.Errorf("failed to get the list of greetings: %w", greetingReadErr)
	}

	if slices.Contains(greetings, newGreeting) {
		return fmt.Errorf("%w: %v", ErrGreetingAlreadyExists, newGreeting)
	}

	return greetingWriter.AddGreeting(ctx, newGreeting)
}

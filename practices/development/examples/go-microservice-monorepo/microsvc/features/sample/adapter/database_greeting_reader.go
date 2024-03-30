package adapter

import (
	"context"
	"fmt"
	"math/rand"

	"example.com/sample/commonlib/database"
)

// DatabaseGreetingReader implements sample.GreetingReader using a live database connection
type DatabaseGreetingReader struct{}

// List implements sample.GreetingReader for DatabaseGreetingReader
func (DatabaseGreetingReader) List(ctx context.Context) ([]string, error) {
	db := database.RetrieveFromContext(ctx)
	var results []string
	selectErr := db.Select(&results, `
		select greetingText from greetings;
	`)
	if selectErr != nil {
		return nil, fmt.Errorf("could not get list of greetings: %w", selectErr)
	}

	return results, nil
}

// RandomGreeting implements sample.GreetingReader for DatabaseGreetingReader
func (rdr DatabaseGreetingReader) RandomGreeting(ctx context.Context) (string, error) {
	greetings, retrieveErr := rdr.List(ctx)
	if retrieveErr != nil {
		return "", retrieveErr
	}

	randomID := rand.Intn(len(greetings))
	return greetings[randomID], nil
}

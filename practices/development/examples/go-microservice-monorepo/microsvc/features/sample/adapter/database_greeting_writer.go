package adapter

import (
	"context"
	"example.com/sample/commonlib/database"
	"fmt"
)

// DatabaseGreetingWriter implements sample.GreetingWriter using a live database connection
type DatabaseGreetingWriter struct{}

// AddGreeting implements GreetingWriter for DatabaseGreetingWriter
func (DatabaseGreetingWriter) AddGreeting(ctx context.Context, newGreeting string) error {
	db := database.RetrieveFromContext(ctx)
	_, insertErr := db.Exec(`
		insert into greetings(greetingText) values (?)
	`, newGreeting)

	if insertErr != nil {
		return fmt.Errorf("failed to add greeting \"%v\": %w", newGreeting, insertErr)
	}

	return nil
}

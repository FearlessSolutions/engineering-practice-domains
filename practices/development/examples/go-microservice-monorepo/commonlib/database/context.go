package database

import (
	"context"
	"database/sql"
	"fmt"

	"example.com/sample/commonlib/request/testhelper"
	"github.com/jmoiron/sqlx"
)

// ctxConnectionKey is where the main database connection is stored in a database context
type ctxConnectionKey struct{}

// ctxTransactionKey is where a transaction handle is stored in a database context
type ctxTransactionKey struct{}

//go:generate mockgen -destination ./context_mocks.go -package database . Connection

// Connection abstracts the API surface of sqlx.DB and sqlx.Tx so database functions can have a single implementation
// regardless of whether a transaction is active or not.
type Connection interface {
	// Get fetches a single row from the database and serializes it into the data structure pointed to by dest.
	Get(dest any, query string, args ...any) error
	// GetContext fetches a single row from tho database and serializes it into the data structure pointed to by dest.
	// Because it accepts a context, it can be cancelled mid-flight, so this is useful for long-running queries
	GetContext(ctx context.Context, dest any, query string, args ...any) error
	// Select fetches many rows into a slice pointed to by dest. Warning: this function pulls the entire result set
	// into memory at once. If you want to break it up, consider using Queryx instead.
	Select(dest any, query string, args ...any) error
	// SelectContext fetches many rows into a slice pointed to by dest. Because it accepts a context, it can be cancelled
	// mid-flight, so this is useful for long-running queries.
	SelectContext(ctx context.Context, dest any, query string, args ...any) error

	// Exec executes the passed query without unmarshalling any data received from the database.
	Exec(query string, args ...any) (sql.Result, error)
	// ExecContext executes the passed query without unmarshalling any data received from the database. It accepts a
	// context, so it can be cancelled mid-flight, making this function good for long-running queries.
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	// NamedExec executes the passed query, binding parameters by name from the passed struct rather than positionally.
	NamedExec(query string, arg any) (sql.Result, error)
	// NamedExecContext executes the passed query, binding parameters by name from the passed struct rather than positionally.
	// It accepts a context so it can be cancelled mid-flight, making this good for long-running queries.
	NamedExecContext(ctx context.Context, query string, arg any) (sql.Result, error)

	// Preparex creates a prepared statement that can be executed multiple times with different data.
	Preparex(query string) (*sqlx.Stmt, error)
	// PreparexContext creates a prepared statement that can be executed multiple times with different data. It can be
	// cancelled mid-flight via the passed context, making it good for functions which use the prepared query for a long
	// time.
	PreparexContext(ctx context.Context, query string) (*sqlx.Stmt, error)
	// PrepareNamed creates a prepared statement using named bind parameters rather than positional ones. The prepared
	// statement can be executed multiple times with different inputs.
	PrepareNamed(query string) (*sqlx.NamedStmt, error)
	// PrepareNamedContext creates a prepared statement using named bind parameters rather than positional ones. The prepared
	// statement can be executed multiple times with different inputs. The statement can be cancelled mid-flight via the
	// passed context, making it good for functions which use the prepared query for a long time.
	PrepareNamedContext(ctx context.Context, query string) (*sqlx.NamedStmt, error)

	// Queryx runs a SQL query and returns an iterator over the retrieved rows which can be serialized into structs as you go.
	Queryx(query string, args ...any) (*sqlx.Rows, error)
	// QueryxContext runs a SQL query and returns an iterator over the retrieved rows which can be serialized into structs as you go.
	// It accepts a context, so it can be cancelled mid-flight, making this good for long queries
	QueryxContext(ctx context.Context, query string, args ...any) (*sqlx.Rows, error)
	// NamedQuery runs a SQL query, binding parameters with named bind parameters rather than positional ones. It returns an iterator
	// over the retrieved rows which can be serialized into structs as you go.
	NamedQuery(query string, arg any) (*sqlx.Rows, error)

	// NOTE: sqlx.Tx does not implement NamedQueryContext. DO NOT ADD THAT FUNCTION TO THIS INTERFACE, IT WILL BREAK THE
	// IMPLICIT IMPLEMENTATION OF CONNECTION FOR sqlx.Tx. Use the database.NamedQueryContext function for this instead.

	// Rebind is used in conjunction with sqlx.In to ensure bind parameters are correct for the current database driver
	Rebind(query string) string
}

// NamedQueryContext runs a SQL query, binding parameters with named bind parameters rather than positional ones. It returns an
// iterator over the retrieved rows which can be serialized into structs as you go. This function accepts a context so it
// can be cancelled mid-flight, making this function good for long-running queries.
func NamedQueryContext(ctx context.Context, cxn Connection, query string, arg any) (*sqlx.Rows, error) {
	switch rawCxn := cxn.(type) {
	case *sqlx.DB:
		return sqlx.NamedQueryContext(ctx, rawCxn, query, arg)
	case *sqlx.Tx:
		return sqlx.NamedQueryContext(ctx, rawCxn, query, arg)
	default:
		panic("NamedQueryContext only accepts database.Connection implementations of *sqlx.DB and *sqlx.Tx! Someone passed a different type!")
	}
}

// CreateDerivativeContext derives a database context from another context. This means creating a child context
// with the database embedded inside so that it can later be extracted with RetrieveFromContext.
func CreateDerivativeContext(ctx context.Context, db *sqlx.DB) context.Context {
	// Don't overwrite the DB if it's already present in the context
	if ctx.Value(ctxConnectionKey{}) != nil {
		return ctx
	}

	return context.WithValue(ctx, ctxConnectionKey{}, db)
}

// CreateDerivativeMockContext derives a database context from another context, attaching a mock connection rather than
// a real database connection.
func CreateDerivativeMockContext(ctx context.Context, mockConnection *MockConnection) context.Context {
	// Don't overwrite the DB if it's already present in the context
	if ctx.Value(ctxConnectionKey{}) != nil {
		return ctx
	}

	return context.WithValue(ctx, ctxConnectionKey{}, mockConnection)
}

// RetrieveFromContext extracts the database connection from the current context. It is expected that some mechanism
// such as middleware.DatabaseContextMiddleware has already added the database to the context via CreateDerivativeContext.
// If the database is not present in the context this function will panic.
func RetrieveFromContext(ctx context.Context) Connection {
	if txConnection := ctx.Value(ctxTransactionKey{}); txConnection != nil {
		return txConnection.(*sqlx.Tx)
	}
	if dbConnection := ctx.Value(ctxConnectionKey{}); dbConnection != nil {
		switch actualConn := dbConnection.(type) {
		case *sqlx.DB:
			return actualConn
		case *MockConnection:
			return actualConn
		}
	}

	panic("Database connection was not present in context! Make sure the database context middleware is installed.")
}

// preparedTransactionContext contains information about a newly created transaction context
type preparedTransactionContext struct {
	transaction         *sqlx.Tx
	passedContext       context.Context
	isNestedTransaction bool
	isMockContext       bool
}

// prepareTransactionContext evaluates the current context to see if a transaction has already been started,
// starting a new one if there's not already an active transaction. This allows transaction functions to be
// harmlessly reentrant
func prepareTransactionContext(parentCtx context.Context) (preparedTransactionContext, error) {
	var preparedTxContext preparedTransactionContext
	// If we're in a mock context (i.e. testing a controller) setting up transactions is a no-op
	if testhelper.IsMockContext(parentCtx) {
		preparedTxContext.passedContext = parentCtx
		preparedTxContext.isMockContext = true
		return preparedTxContext, nil
	}

	switch rawCxn := RetrieveFromContext(parentCtx).(type) {
	case *sqlx.Tx:
		preparedTxContext.transaction = rawCxn
		preparedTxContext.isNestedTransaction = true
		preparedTxContext.passedContext = parentCtx
	case *sqlx.DB:
		newTx, txBeginErr := rawCxn.Beginx()
		if txBeginErr != nil {
			return preparedTxContext, txBeginErr
		}

		preparedTxContext.transaction = newTx
		preparedTxContext.passedContext = context.WithValue(parentCtx, ctxTransactionKey{}, newTx)
	}

	return preparedTxContext, nil
}

// rollbackOnFailureOrCommit runs transaction finalization logic if the prepared transaction context is the one which
// initially started the transaction. It evaluates the error returned from the nested operation, committing the changes
// if no error was returned or rolling back if there was an error. If there is an error during rollback, the original error
// is wrapped in such a way that it will still be accessible via errors.Is and errors.As.
func rollbackOnFailureOrCommit(operationError error, preparedCtx preparedTransactionContext) error {
	// Do nothing in a mock context (i.e. testing a controller), transaction setup/teardown is a no-op in that case
	if preparedCtx.isMockContext {
		return operationError
	}

	returnedError := operationError
	if !preparedCtx.isNestedTransaction {
		if operationError != nil {
			rollbackErr := preparedCtx.transaction.Rollback()
			if rollbackErr != nil {
				returnedError = fmt.Errorf("rollback failed when operation returned an error (%w): %w", operationError, rollbackErr)
			}
		} else {
			returnedError = preparedCtx.transaction.Commit()
		}
	}
	return returnedError
}

// WithTransaction initiates a database transaction which finalizes after the passed function is executed. By "finalizes", this
// means the transaction is committed if the passed function does not return an error, otherwise the transaction is safely rolled back.
// This function is safely reentrant, meaning calling WithTransaction inside a WithTransaction function won't do anything.
//
// The error returned from the passed function will be returned from this function. This function may also return
// database errors associated with starting or finalizing the transaction. If this occurs on a rollback, the original
// error will be wrapped and accessible via errors.Is or errors.As.
func WithTransaction(ctx context.Context, operation func(ctx context.Context) error) error {
	preparedCtx, prepareErr := prepareTransactionContext(ctx)
	if prepareErr != nil {
		return prepareErr
	}

	operationError := operation(preparedCtx.passedContext)
	return rollbackOnFailureOrCommit(operationError, preparedCtx)
}

// WithTransactionReturning initiates a database transaction which finalizes after the passed function is executed. By "finalizes", this
// means the transaction is committed if the passed function does not return an error, otherwise the transaction is safely rolled back.
// This function is safely reentrant, meaning calling WithTransactionReturning inside a WithTransactionReturning function won't do anything.
//
// The error returned from the passed function will be returned from this function. This function may also return
// database errors associated with starting or finalizing the transaction. If this occurs on a rollback, the original
// error will be wrapped and accessible via errors.Is or errors.As.
func WithTransactionReturning[ReturnValue any](ctx context.Context, operation func(ctx context.Context) (ReturnValue, error)) (ReturnValue, error) {
	preparedCtx, prepareErr := prepareTransactionContext(ctx)
	if prepareErr != nil {
		var zeroValue ReturnValue
		return zeroValue, prepareErr
	}

	returnValue, operationErr := operation(preparedCtx.passedContext)
	return returnValue, rollbackOnFailureOrCommit(operationErr, preparedCtx)
}

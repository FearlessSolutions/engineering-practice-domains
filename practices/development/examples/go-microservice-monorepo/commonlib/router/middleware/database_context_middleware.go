package middleware

import (
	"example.com/sample/commonlib/database"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

// DatabaseContextMiddleware attaches a database connection to the HTTP request context which can be extracted in
// driven ports. Exposing the database connection in the request context also allows context-wide database transactions
// to be performed without exposing it in the driven port interface.
func DatabaseContextMiddleware(db *sqlx.DB) echo.MiddlewareFunc {
	// This is the actual middleware function that gets called for every request
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		// This function contains the middleware logic which runs on a single request
		return func(ctx echo.Context) error {
			request := ctx.Request()
			reqCtx := request.Context()

			dbConnectionCtx := database.CreateDerivativeContext(reqCtx, db)
			requestWithDBContext := request.WithContext(dbConnectionCtx)
			ctx.SetRequest(requestWithDBContext)

			return next(ctx)
		}
	}
}

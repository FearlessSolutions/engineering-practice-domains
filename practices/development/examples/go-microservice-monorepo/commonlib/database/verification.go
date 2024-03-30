package database

import (
	"errors"
	"time"

	"example.com/sample/commonlib/logger"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// MustBeConnected asserts that the database is actually connected. It attempts to ping the database once every 10 seconds
// for 5 minutes, and if a connection is never successfully established this function panics and kills the server.
func MustBeConnected(db *sqlx.DB) {
	startTime := time.Now()
	for {
		// Give up after 5 minutes
		if time.Now().After(startTime.Add(5 * time.Minute)) {
			logger.Log.Fatal("Database could not be reached after 5 minutes! Shutting down.")
		}
		if connectErr := db.Ping(); connectErr != nil {
			// Make sure it's not something like a "password is wrong" error or something
			var target *mysql.MySQLError
			if errors.As(connectErr, &target) {
				logger.Log.Fatal("Database successfully connected, but connection options may be wrong. Shutting down.", zap.Error(connectErr))
			}

			// If the problem wasn't something a successful connection told us, retry connecting in 10 seconds
			logger.Log.Warn("Could not reach the database, retrying in 10s...", zap.Error(connectErr))
			time.Sleep(10 * time.Second)
		} else {
			break
		}
	}
}

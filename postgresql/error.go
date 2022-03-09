package postgresql

import (
	"errors"
)

var (
	// ErrEnvironmentVariable is returned when the worker fails to find
	// the environment variables necessary to connect to the database.
	ErrEnvironmentVariable = errors.New("could not find all necessary environment variables")
	// ErrDatabaseConnection is returned when the worker fails to connect to
	// the database.
	ErrDatabaseConnection = errors.New("database connection error")
	// ErrDatabasePing is returned when the worker fails to ping the
	// database.
	ErrDatabasePing = errors.New("database ping error")
	// ErrPreparingStmt is returned when the worker fails to prepare
	// a prepared statement.
	ErrPreparingStmt = errors.New("error executing prepared statement")
	// ErrExecutingPreparedStmt is returned when the worker fails to execute
	// a query with a prepared statement.
	ErrExecutingPreparedStmt = errors.New("error executing prepared statement")
	// ErrScanningRows is returned when the worker fails to scan the rows
	// returned by a query.
	ErrScanningRows = errors.New("error trying to scan result rows")
	// ErrExecutingRollback is returned when the worker fails to rollback
	// in a transaction.
	ErrExecutingRollback = errors.New("error trying to rollback a transaction in database")
)

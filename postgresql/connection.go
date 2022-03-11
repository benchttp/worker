package postgresql

import (
	"database/sql"
	"fmt"

	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres" // blank import

	"github.com/benchttp/worker/benchttp"
)

// InsertionService implements benchttp.InsertionService interface.
type InsertionService struct {
	db *sql.DB
}

// NewInsertionService connects to the database and provides an InsertionService
// implementing benchttp.InsertionService, which provides a method to insert
// data in database.
func NewInsertionService(config benchttp.Config) (InsertionService, error) {
	dbURI := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host,
		config.User,
		config.Password,
		config.DBName)

	db, err := sql.Open("cloudsqlpostgres", dbURI)
	if err != nil {
		return InsertionService{}, err
	}

	err = db.Ping()
	if err != nil {
		return InsertionService{}, err
	}

	db.SetMaxIdleConns(config.IdleConn)
	db.SetMaxOpenConns(config.MaxConn)

	return InsertionService{db}, nil
}

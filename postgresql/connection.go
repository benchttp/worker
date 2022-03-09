package postgresql

import (
	"database/sql"
	"fmt"

	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres" // blank import

	"github.com/benchttp/worker"
)

type InsertionService struct {
	db *sql.DB
}

func NewInsertionService(config worker.Config) (InsertionService, error) {
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

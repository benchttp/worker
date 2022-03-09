package postgresql

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres" // blank import
)

type StatsService struct {
	db *sql.DB
}

func NewStatsService(config Config) (StatsService, error) {
	dbURI := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host,
		config.User,
		config.Password,
		config.DBName)

	db, err := sql.Open("cloudsqlpostgres", dbURI)
	if err != nil {
		return StatsService{}, err
	}

	err = db.Ping()
	if err != nil {
		return StatsService{}, err
	}

	db.SetMaxIdleConns(config.IdleConn)
	db.SetMaxOpenConns(config.MaxConn)

	return StatsService{db}, nil
}

type Config struct {
	Host     string
	User     string
	Password string
	DBName   string
	IdleConn int
	MaxConn  int
}

func GetConfigFromEnvVariables() (Config, error) {
	var config Config

	config.Host = os.Getenv("PSQL_HOST")
	if config.Host == "" {
		return config, errors.New("PSQL_HOST environment variable not found")
	}

	config.User = os.Getenv("PSQL_USER")
	if config.User == "" {
		return config, errors.New("PSQL_USER environment variable not found")
	}

	config.Password = os.Getenv("PSQL_PASSWORD")
	if config.Password == "" {
		return config, errors.New("PSQL_PASSWORD environment variable not found")
	}

	config.DBName = os.Getenv("PSQL_NAME")
	if config.DBName == "" {
		return config, errors.New("PSQL_NAME environment variable not found")
	}

	config.IdleConn = 10
	config.MaxConn = 25

	return config, nil
}

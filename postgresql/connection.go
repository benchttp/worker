package postgresql

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"

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
		return StatsService{}, ErrDatabaseConnection
	}

	err = db.Ping()
	if err != nil {
		return StatsService{}, ErrDatabasePing
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

	err := godotenv.Load(".env")
	// no error returned here because .env is not deployed
	if err != nil {
		fmt.Println("no .env file found")
	}

	config.Host = os.Getenv("PSQL_HOST")
	if config.Host == "" {
		return config, ErrEnvironmentVariable
	}

	config.User = os.Getenv("PSQL_USER")
	if config.User == "" {
		return config, ErrEnvironmentVariable
	}

	config.Password = os.Getenv("PSQL_PASSWORD")
	if config.Password == "" {
		return config, ErrEnvironmentVariable
	}

	config.DBName = os.Getenv("PSQL_NAME")
	if config.DBName == "" {
		return config, ErrEnvironmentVariable
	}

	config.IdleConn = 10
	config.MaxConn = 25

	return config, nil
}

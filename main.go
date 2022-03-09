package worker

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/googleapis/google-cloudevents-go/cloud/firestore/v1"

	"github.com/benchttp/worker/firestoreconv"
	"github.com/benchttp/worker/postgresql"
	"github.com/benchttp/worker/stats"
)

// Digest is a Cloud Function triggered by a Firestore create document
// event to extract, compute and store statistics of a Benchttp run.
func Digest(ctx context.Context, e firestore.DocumentEventData) error {
	r, err := firestoreconv.Report(e.Value)
	if err != nil {
		return err
	}

	cfg, err := envConfig()
	insertionService, err := postgresql.NewInsertionService(cfg)

	codes, times := r.Benchmark.Values()

	timestats, err := stats.ComputeCommon(times)
	if err != nil {
		return err
	}

	log.Printf("timestats: %s", timestats.StringAsTime())

	codestats, err := stats.ComputeStatusDistribution(codes)
	if err != nil {
		return err
	}

	log.Printf("codestats: %+v", codestats)

	return nil
}

type Config struct {
	Host     string
	User     string
	Password string
	DBName   string
	IdleConn int
	MaxConn  int
}

func envConfig() (Config, error) {
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

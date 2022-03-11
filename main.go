package worker

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"github.com/googleapis/google-cloudevents-go/cloud/firestore/v1"

	"github.com/benchttp/worker/benchttp"
	"github.com/benchttp/worker/firestoreconv"
	"github.com/benchttp/worker/postgresql"
	"github.com/benchttp/worker/stats"
)

// Digest is a Cloud Function triggered by a Firestore create document
// event to extract, compute and store statistics of a Benchttp run.
// It also stores the computed data in a SQL database.
func Digest(ctx context.Context, e firestore.DocumentEventData) error {
	r, err := firestoreconv.Report(e.Value)
	if err != nil {
		return err
	}

	cfg, err := envConfig()
	if err != nil {
		return err
	}
	insertionService, err := postgresql.NewInsertionService(cfg)
	if err != nil {
		return err
	}

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

	// TO DO: get user id. Using "1" here for the moment.
	statsToInsert := buildStats(timestats, codestats, "firestore_id", "1")

	if err := insertionService.Insert(statsToInsert); err != nil {
		return err
	}

	return nil
}

func envConfig() (benchttp.Config, error) {
	var config benchttp.Config

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

// buildStats builds a benchttp.Stats object.
// Descriptor.FinishedAt is set at time.now().
func buildStats(timestats stats.Common, codestats stats.StatusDistribution, reportID, userID string) benchttp.Stats {
	computedstats := benchttp.Stats{
		Descriptor: benchttp.StatsDescriptor{
			ID:         reportID,
			UserID:     userID,
			FinishedAt: time.Now(),
		},
		Time: timestats,
		Code: codestats,
	}
	return computedstats
}

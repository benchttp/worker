package worker

import (
	"context"
	"log"

	"github.com/googleapis/google-cloudevents-go/cloud/firestore/v1"

	"github.com/benchttp/worker/firestoreconv"
	"github.com/benchttp/worker/stats"
)

// Digest is a Cloud Function triggered by a Firestore create document
// event to extract, compute and store statistics of a Benchttp run.
func Digest(ctx context.Context, e firestore.DocumentEventData) error {
	r, err := firestoreconv.Report(e.Value)
	if err != nil {
		return err
	}

	data := r.Benchmark.Times()

	s, err := stats.Compute(data)
	if err != nil {
		return err
	}

	log.Print(s.String())

	return nil
}

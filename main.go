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
	log.Printf("→ firestore protobuf document: %v", e)

	r, err := firestoreconv.Report(e.Value)
	if err != nil {
		return err
	}

	data := r.Benchmark.Times()

	s, err := stats.Compute(data)
	if err != nil {
		return err
	}

	b, err := s.MarshalJSON()
	if err != nil {
		log.Printf("warning: failed to marshal stats: %s", err)
		return nil
	}
	log.Print(string(b))

	return nil
}

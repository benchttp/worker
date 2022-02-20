package worker

import (
	"context"
	"log"

	"github.com/googleapis/google-cloudevents-go/cloud/firestore/v1"

	"github.com/benchttp/worker/firestoreconv"
	"github.com/benchttp/worker/internal"
)

// Digest is a Cloud Function triggered by a Firestore create document
// event to extract, compute and store statistics of a Benchttp run.
func Digest(ctx context.Context, e firestore.DocumentEventData) error {
	log.Printf("→ firestore protobuf document: %v", e)

	b, err := firestoreconv.ToBenchmark(e.Value)
	if err != nil {
		return err
	}

	stats, err := internal.ComputeStats(b.Report)
	if err != nil {
		return err
	}

	log.Printf("→ computed stats: %v", stats)

	return nil
}

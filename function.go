package worker

import (
	"context"
	"log"

	"github.com/googleapis/google-cloudevents-go/cloud/firestore/v1"

	"github.com/benchttp/worker/firestoreconv"
	"github.com/benchttp/worker/internal"
)

func Handle(ctx context.Context, e firestore.DocumentEventData) error {
	log.Printf("event: %v", e)

	report, err := firestoreconv.ToReport(e.Value)
	if err != nil {
		return err
	}

	log.Printf("report: %v", report)

	stats, err := internal.ComputeStats(report)
	if err != nil {
		return err
	}

	log.Printf("stats: %v", stats)

	return nil
}

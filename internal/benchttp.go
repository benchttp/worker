package internal

import (
	"time"

	"github.com/benchttp/worker/timestats"
)

// Benchmark represents the result of a Benchttp benchmark run.
type Benchmark struct {
	Report Report
}

// Report represents the detailed collection of requests done
// during a Benchttp benchmark run.
type Report struct {
	Records []Record
}

// Record represents the summary of a HTTP response.
type Record struct {
	Time time.Duration
}

func timeData(rec []Record) []float64 {
	floats := make([]float64, len(rec))
	for i, v := range rec {
		floats[i] = float64(v.Time)
	}
	return floats
}

// ComputeStats computes the statistics of a Report.
func ComputeStats(r Report) (timestats.Stats, error) {
	data := timeData(r.Records)
	stats, err := timestats.Compute(data)
	if err != nil {
		return timestats.Stats{}, err
	}
	return stats, nil
}

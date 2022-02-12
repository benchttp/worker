package internal

import (
	"time"

	"github.com/benchttp/worker/timestats"
)

type Report struct {
	Records []Record
}

type Record struct {
	Time time.Duration
}

func (r Report) timeData() []float64 {
	floats := make([]float64, len(r.Records))
	for i, v := range r.Records {
		floats[i] = float64(v.Time)
	}
	return floats
}

// ComputeStats computes the statistics of a Report.
func ComputeStats(report Report) (timestats.Stats, error) {
	data := report.timeData()
	stats, err := timestats.Compute(data)
	if err != nil {
		return timestats.Stats{}, err
	}
	return stats, nil
}

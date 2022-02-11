// Package timestats provides a wrapper around github.com/montanaflynn/stats
// to compute common statistics for a dataset whose values are interpreted
// as time.Duration.
//
// The package requires the dataset to be a slice of float64. The caller has
// the responsibility of transforming the data into a usable dataset for
// timestats.Compute.
//
// For example:
//	type MyType struct {
//		times time.Duration[]
//	}
//
//	func (t MyType) FloatSlice() []float64 {
//		floats := make([]float64, len(t.times))
//		for i, v := range t.times {
//			floats[i] = float64(v)
//		}
//		return floats
//	}
//
//	var t MyType
//
//	stats, err := timestats.Compute(t.FloatSlice())
//
package timestats

import (
	"fmt"
	"time"

	"github.com/montanaflynn/stats"
)

var deciles = [9]float64{10, 20, 30, 40, 50, 60, 70, 80, 90}

// Stats represents the statistics computed from a given dataset.
type Stats struct {
	Min      time.Duration `json:"min"`
	Max      time.Duration `json:"max"`
	Mean     time.Duration `json:"mean"`
	Median   time.Duration `json:"median"`
	Variance time.Duration `json:"variance"`
	// Deciles is an array composed of the nine deciles.
	// Deciles[0] is the first decile (10%), Deciles[8]
	// is the 9th decile (90%).
	Deciles [9]time.Duration `json:"deciles"`
}

// Compute computes the common statistics for a dataset.
// ErrCompute if any error occurs during the computation.
// When returning an error, the value of Stats may be
// partially written.
func Compute(data []float64) (Stats, error) {
	if len(data) == 0 {
		return Stats{}, ErrEmptySlice
	}

	output := Stats{}
	issues := []string{}

	// Handle flat statistics.
	for _, v := range []struct {
		name string
		f    func(input stats.Float64Data) (float64, error)
		dst  *time.Duration
	}{
		{"min", stats.Min, &output.Min},
		{"max", stats.Max, &output.Max},
		{"mean", stats.Mean, &output.Mean},
		{"median", stats.Median, &output.Median},
		{"variance", stats.Variance, &output.Variance},
	} {
		stat, err := computeStat(v.f, data, v.name)
		*v.dst = time.Duration(stat)
		if err != nil {
			issues = append(issues, err.Error())
		}
	}

	// Handle exception case Stats.Deciles.
	for i, percent := range deciles {
		stat, err := computeDecile(percent, data)
		output.Deciles[i] = time.Duration(stat)
		if err != nil {
			issues = append(issues, err.Error())
		}
	}

	if len(issues) > 0 {
		return output, &ComputeError{issues}
	}
	return output, nil
}

func computeStat(f func(in stats.Float64Data) (float64, error), data []float64, name string) (float64, error) {
	stat, err := f(data)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", name, err)
	}
	return stat, nil
}

func computeDecile(percent float64, data []float64) (float64, error) {
	decile, err := stats.Percentile(data, percent)
	if err != nil {
		name := "percentile"
		return 0, fmt.Errorf("%s: computing %s percentile with slice of length %d: %w",
			name, ordinal(int(percent/100+1)), len(data), err)
	}
	return decile, nil
}

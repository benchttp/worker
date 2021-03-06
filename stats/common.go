// Package stats provides a wrapper around github.com/montanaflynn/stats
// to compute common statistics for a dataset.
//
// The package requires the dataset to be a slice of float64. The caller has
// the responsibility of transforming the data into a usable dataset for
// stats.Compute.
//
package stats

import (
	"fmt"
	"time"

	"github.com/montanaflynn/stats"
)

// Common represents the common statistics computed from a given dataset.
type Common struct {
	Min    float64
	Max    float64
	Mean   float64
	Median float64
	StdDev float64

	// Deciles is an array composed of the nine deciles.
	// Deciles[0] is the first decile (10%), Deciles[8]
	// is the 9th decile (90%).
	Deciles [9]float64
}

// ComputeCommon computes the common statistics for a dataset. Returns ErrCompute
// if any error occurs during the computation. When returning an error, the
// value of Stats struct may be partially written.
func ComputeCommon(data []float64) (Common, error) {
	if len(data) == 0 {
		return Common{}, ErrEmptySlice
	}

	out := Common{}
	errs := []string{}

	out.Max, errs = pipe("max", errs)(stats.Max(data))
	out.Min, errs = pipe("min", errs)(stats.Min(data))
	out.Mean, errs = pipe("mean", errs)(stats.Mean(data))
	out.Median, errs = pipe("median", errs)(stats.Median(data))
	out.StdDev, errs = pipe("stddev", errs)(stats.StandardDeviation(data))

	for i, p := range deciles {
		n := fmt.Sprintf("%s decile", ordinal(i+1))
		out.Deciles[i], errs = pipe(n, errs)(stats.Percentile(data, p))
	}

	if len(errs) > 0 {
		return out, &ComputeError{errs}
	}

	return out, nil
}

var deciles = [9]float64{10, 20, 30, 40, 50, 60, 70, 80, 90}

func pipe(name string, errs []string) func(float64, error) (float64, []string) {
	return func(stat float64, err error) (float64, []string) {
		if err != nil {
			errs = append(errs, fmt.Sprintf("computing %s: %s", name, err))
		}
		return stat, errs
	}
}

// StringAsTime returns a string representing the statistics
// expressed in time.Duration.
//	{min:100ms max:200ms ...}
func (s *Common) StringAsTime() string {
	v := struct {
		Min     string
		Max     string
		Mean    string
		Median  string
		StdDev  string
		Deciles [9]string
	}{
		Min:    time.Duration(s.Min).String(),
		Max:    time.Duration(s.Max).String(),
		Mean:   time.Duration(s.Mean).String(),
		Median: time.Duration(s.Median).String(),
		StdDev: time.Duration(s.StdDev).String(),
	}
	for i, d := range s.Deciles {
		v.Deciles[i] = time.Duration(d).String()
	}
	return fmt.Sprintf("%+v", v)
}

package timestats

import (
	"time"

	"github.com/montanaflynn/stats"
)

var deciles = [9]float64{10, 20, 30, 40, 50, 60, 70, 80, 90}

type Stats struct {
	Min      time.Duration    `json:"min"`
	Max      time.Duration    `json:"max"`
	Mean     time.Duration    `json:"mean"`
	Median   time.Duration    `json:"median"`
	Variance time.Duration    `json:"variance"`
	Deciles  [9]time.Duration `json:"deciles"`
}

func Compute(data stats.Float64Data) Stats {
	var s Stats

	s.Min = silent(data, stats.Min)
	s.Max = silent(data, stats.Max)
	s.Mean = silent(data, stats.Mean)
	s.Median = silent(data, stats.Median)
	s.Variance = silent(data, stats.Variance)

	for i, percent := range deciles {
		decile, _ := stats.Percentile(data, percent)
		s.Deciles[i] = time.Duration(decile)
	}
	return s
}

func silent(data stats.Float64Data, f func(input stats.Float64Data) (float64, error)) time.Duration {
	stat, err := f(data)
	if err != nil {
		panic(err) // TODO
	}
	return time.Duration(stat)
}

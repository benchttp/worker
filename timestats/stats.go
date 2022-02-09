// Package timestats provides a wrapper around github.com/montanaflynn/stats
// to compute common statistics for a dataset of time.Duration values.
//
// Any type implementing timestats.FloatSlicer or slice of timestats.Floater
// can be transformed to a dataset compatible with the package.
//
// For example, using timestats.FloatSlicer:
//	type MySlice []int64
//	func (s MySlice) FloatSlice() []float64 {...}
//
//	dataset := timestats.Transform(MySlice{1, 2, 3})
//	stats, err := timestats.Compute(dataset)
//
// Using timestats.Floater:
//	type MyType int64
//	func (t MyType) Float() float64 {...}
//
//	dataset, err := timestats.TransformIter([]MyType{1, 2, 3})
//	stats, err := timestats.Compute(dataset)
//
package timestats

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
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
func Compute(data stats.Float64Data) (Stats, error) {
	if len(data) == 0 {
		return Stats{}, ErrEmptySlice
	}

	output := Stats{}
	issues := []string{}

	// Handle flat statistics.
	for _, v := range []struct {
		f   func(input stats.Float64Data) (float64, error)
		dst *time.Duration
	}{
		{stats.Min, &output.Min},
		{stats.Max, &output.Max},
		{stats.Mean, &output.Mean},
		{stats.Median, &output.Median},
		{stats.Variance, &output.Variance},
	} {
		*v.dst, issues = computeStat(v.f, data, issues)
	}

	// Handle exception case Stats.Deciles.
	for i, percent := range deciles {
		output.Deciles[i], issues = computeDecile(percent, data, issues)
	}

	if len(issues) > 0 {
		return output, &ErrCompute{issues}
	}
	return output, nil
}

func computeStat(f func(in stats.Float64Data) (float64, error), data stats.Float64Data, e []string) (time.Duration, []string) {
	stat, err := f(data)
	if err != nil {
		funcname := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
		// github.com/montanaflynn/stats.Min -> Min
		funcname = strings.Split(funcname, ".")[2:][0]
		e = append(e, fmt.Sprintf("%s: %s", funcname, err))
		return 0, e
	}
	return time.Duration(stat), e
}

func computeDecile(percent float64, data stats.Float64Data, e []string) (time.Duration, []string) {
	decile, err := stats.Percentile(data, percent)
	if err != nil {
		funcname := "Percentile"
		e = append(e, fmt.Sprintf("%s: %s", funcname, err))
		return 0, e
	}
	return time.Duration(decile), e
}

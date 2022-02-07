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

type Stats struct {
	Min      time.Duration    `json:"min"`
	Max      time.Duration    `json:"max"`
	Mean     time.Duration    `json:"mean"`
	Median   time.Duration    `json:"median"`
	Variance time.Duration    `json:"variance"`
	Deciles  [9]time.Duration `json:"deciles"`
}

func Compute(data stats.Float64Data) (Stats, error) {
	s := Stats{}
	issues := []string{}

	// Handle flat statistics.
	for _, v := range []struct {
		f   func(input stats.Float64Data) (float64, error)
		dst *time.Duration
	}{
		{stats.Min, &s.Min},
		{stats.Max, &s.Max},
		{stats.Mean, &s.Mean},
		{stats.Median, &s.Median},
		{stats.Variance, &s.Variance},
	} {
		*v.dst, issues = computeStat(v.f, data, issues)
	}

	// Handle exception case Stats.Deciles.
	for i, percent := range deciles {
		s.Deciles[i], issues = computeDecile(percent, data, issues)
	}

	if len(issues) > 0 {
		return s, &ErrCompute{issues}
	}
	return s, nil
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
	}
	return time.Duration(decile), e
}

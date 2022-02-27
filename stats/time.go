package stats

import (
	"fmt"
	"time"
)

// String returns a string representing the statistics
// expressed in time.Duration.
//	{min:100ms max:200ms ...}
func (s *Stats) String() string {
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

package stats

import (
	"encoding/json"
	"time"
)

func (s *Stats) MarshalJSON() ([]byte, error) {
	v := &struct {
		Min     time.Duration    `json:"min"`
		Max     time.Duration    `json:"max"`
		Mean    time.Duration    `json:"mean"`
		Median  time.Duration    `json:"median"`
		StdDev  time.Duration    `json:"stddev"`
		Deciles [9]time.Duration `json:"deciles"`
	}{
		Min:    time.Duration(s.Min),
		Max:    time.Duration(s.Max),
		Mean:   time.Duration(s.Mean),
		Median: time.Duration(s.Median),
		StdDev: time.Duration(s.StdDev),
	}
	return json.MarshalIndent(v, "", "  ")
}

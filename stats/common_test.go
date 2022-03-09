package stats_test

import (
	"errors"
	"math/rand"
	"reflect"
	"testing"

	"github.com/benchttp/worker/stats"
)

func TestComputeCommon(t *testing.T) {
	t.Run("passing invalid dataset returns error", func(t *testing.T) {
		for _, testcase := range []struct {
			name string
			data []float64
			want error
			zero bool // whether or not ComputeCommon response may be partially written
		}{
			{
				name: "empty dataset",
				data: []float64{},
				want: stats.ErrEmptySlice,
				zero: true,
			},
			{
				name: "not enough values",
				data: newDataStub(2), // not enough for 9 deciles
				want: &stats.ComputeError{},
				zero: false,
			},
		} {
			t.Run(testcase.name, func(t *testing.T) {
				res, err := stats.ComputeCommon(testcase.data)

				if err == nil {
					t.Error("want error, got none")
				}

				if !errors.As(err, &testcase.want) {
					t.Errorf("want %T, got %+v", testcase.want, err)
				}

				switch {
				case testcase.zero && !reflect.ValueOf(res).IsZero():
					t.Errorf("want stats output to be zero value, got %+v", res)
				case !testcase.zero && reflect.ValueOf(res).IsZero():
					t.Error("want stats output to be non-zero value, got zero value")
				}
			})
		}
	})

	t.Run("happy path", func(t *testing.T) {
		var (
			size = 10000
			want = stats.Common{
				Min:    1,
				Max:    10000,
				Mean:   5000,
				Median: 5000,
				// StdDev:  0,
				Deciles: [9]float64{1000, 2000, 3000, 4000, 5000, 6000, 7000, 8000, 9000},
			}
		)

		data := newDataStub(size)
		got, err := stats.ComputeCommon(data)
		if err != nil {
			t.Fatalf("want nil error, got %v", err)
		}

		if reflect.ValueOf(got).IsZero() {
			t.Error("want stats output to be non-zero value, got zero value")
		}

		for _, stat := range []struct {
			name string
			want float64
			got  float64
		}{
			{"min", want.Min, got.Min},
			{"max", want.Max, got.Max},
			{"mean", want.Mean, got.Mean},
			{"median", want.Median, got.Median},
			// {"standard deviation", want.StdDev, got.StdDev},
		} {
			if !approxEqual(stat.got, stat.want, 1) {
				t.Errorf("%s: want %f, got %f", stat.name, stat.want, stat.got)
			}
		}

		for i, got := range got.Deciles {
			if !approxEqual(got, want.Deciles[i], 1) {
				t.Errorf("decile %d: want %f, got %f", (i+1)*100, got, want.Deciles[i])
			}
		}
	})
}

// newDataStub returns a slice of float64 of length size. It is filled
// with randomly arranged numbers but with the assurance of containing
// exactly one onccurrence of each number from 0 to size. For example:
//
//	newDataStub(10) -> [1 5 3 7 9 2 4 6 8 0]
//
// The returned dataset offers predictable output when computing common
// statistics such as min, max, mean, etc.
func newDataStub(size int) []float64 {
	floats := make([]float64, size)

	for i, v := range rand.Perm(size) {
		floats[i] = float64(v + 1)
	}
	return floats
}

// approxEqual returns true if val is equal to target with a margin of error.
func approxEqual(val, target, margin float64) bool {
	return val >= target-margin && val <= target+margin
}

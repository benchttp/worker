package timestats_test

import (
	"errors"
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/montanaflynn/stats"

	"github.com/benchttp/worker/timestats"
)

// The following test simply tests that timestats.Compute return a non-zero
// value of timestats.Stats. The actual statistics computation is not tested,
// as it would be equivalent to testing the third party  package itself.
// We can safely assume it is correctly tested.

func TestCompute(t *testing.T) {
	t.Run("passing invalid dataset returns error", func(t *testing.T) {
		for _, testcase := range []struct {
			name string
			data stats.Float64Data
			want interface{}
			zero bool // whether or not Compute response may be partially written
		}{
			{
				name: "empty dataset",
				data: stats.Float64Data{},
				want: timestats.ErrEmptySlice,
				zero: true,
			},
			{
				name: "not enough values",
				data: stats.Float64Data{1, 1, 1, 1, 1, 1, 1, 1, 1}, // 9 values is not enough for 9 deciles
				want: &timestats.ComputeError{},
				zero: false,
			},
		} {
			t.Run(testcase.name, func(t *testing.T) {
				res, err := timestats.Compute(testcase.data)

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
			min  = 1 * time.Nanosecond
			max  = time.Duration(size) * time.Nanosecond
			mean = 5000 * time.Nanosecond
		)

		data := newDataStub(size)
		res, err := timestats.Compute(data)
		if err != nil {
			t.Fatalf("want nil error, got %v", err)
		}

		if reflect.ValueOf(res).IsZero() {
			t.Error("want stats output to be non-zero value, got zero value")
		}

		if res.Min != min {
			t.Errorf("want min as %s, got %s", res.Min, min)
		}
		if res.Max != max {
			t.Errorf("want max as %s, got %s", res.Max, max)
		}
		if res.Mean != mean {
			t.Errorf("want mean as %s, got %s", res.Mean, mean)
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
func newDataStub(size int) stats.Float64Data {
	floats := make([]float64, size)

	for i, v := range rand.Perm(size) {
		floats[i] = float64(v + 1)
	}
	return floats
}

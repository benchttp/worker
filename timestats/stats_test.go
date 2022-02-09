package timestats_test

import (
	"errors"
	"reflect"
	"testing"

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
		data := stats.Float64Data{1, 10, 100, 200, 300, 400, 500, 600, 700, 800, 900, 1000}

		res, err := timestats.Compute(data)
		if err != nil {
			t.Fatalf("want nil error, got %v", err)
		}

		if reflect.ValueOf(res).IsZero() {
			t.Error("want stats output to be non-zero value, got zero value")
		}
	})
}

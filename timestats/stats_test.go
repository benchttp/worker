package timestats_test

import (
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
	t.Run("passing invalid dataset returns ErrCompute", func(t *testing.T) {
		empty := stats.Float64Data{}

		res, err := timestats.Compute(empty)
		if err == nil {
			t.Error("want error, got none")
		}
		if reflect.ValueOf(res).IsZero() {
			t.Errorf("want timestats.Stats to be non-zero value, got %+v", res)
		}
	})

	t.Run("happy path", func(t *testing.T) {
		data := stats.Float64Data{1, 10, 100, 200, 300, 400, 500, 600, 700, 800, 900, 1000}

		res, err := timestats.Compute(data)
		if err != nil {
			t.Fatalf("want nil error, got %v", err)
		}
		if reflect.ValueOf(res).IsZero() {
			t.Errorf("want timestats.Stats to be non-zero value, got %+v", res)
		}
	})
}

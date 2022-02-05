package timestats_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/benchttp/worker/timestats"
)

// TestCompute simply tests that timestats.Compute correctly
// returns a non-zero value timestats.Stats. The actual statistics
// computed are not tested, as we would simply test that the third
// party package is correctly tested.
func TestCompute(t *testing.T) {
	fakes := fakes{
		{time: 1 * time.Nanosecond},
		{time: 10 * time.Nanosecond},
		{time: 100 * time.Nanosecond},
		{time: 200 * time.Nanosecond},
		{time: 300 * time.Nanosecond},
		{time: 400 * time.Nanosecond},
		{time: 500 * time.Nanosecond},
		{time: 600 * time.Nanosecond},
		{time: 700 * time.Nanosecond},
		{time: 800 * time.Nanosecond},
		{time: 900 * time.Nanosecond},
		{time: 1000 * time.Nanosecond},
	}
	raw, err := timestats.Transform(fakes)
	if err != nil {
		t.Fatalf("want nil error, got %v", err)
	}

	s := timestats.Compute(raw)
	if reflect.ValueOf(s).IsZero() {
		t.Errorf("want timestats.Stats to be non-zero value, got %+v", s)
	}
}

package stats_test

import (
	"reflect"
	"testing"

	"github.com/benchttp/worker/stats"
)

func TestComputeStatusDistribution(t *testing.T) {
	t.Run("returns an error for invalid status codes", func(t *testing.T) {
		data := []int{200, 400, 1000}
		want := stats.ComputeError{[]string{"1000 is not a valid HTTP status code"}}

		_, err := stats.ComputeStatusDistribution(data)
		if err == nil {
			t.Error("want error, got none")
		}

		if err.Error() != want.Error() {
			t.Errorf("want %T, got %+v", want, err)
		}
	})

	t.Run("happy path", func(t *testing.T) {
		data := []int{200, 201, 300, 400, 403, 404, 500, 502}
		want := stats.StatusDistribution{
			Status1xx: 0,
			Status2xx: 2,
			Status3xx: 1,
			Status4xx: 3,
			Status5xx: 2,
		}

		got, err := stats.ComputeStatusDistribution(data)
		if err != nil {
			t.Fatalf("want nil error, got %v", err)
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("incorrect computation of status distribution: want %v, got %v", want, got)
		}
	})
}

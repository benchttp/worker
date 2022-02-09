package timestats_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/montanaflynn/stats"

	"github.com/benchttp/worker/timestats"
)

type floater struct {
	time time.Duration
	any  interface{} //nolint
}

func (e floater) Float() float64 {
	return float64(e.time)
}

func TestTransformIter(t *testing.T) {
	for _, testcase := range []struct {
		name    string
		raw     interface{}
		want    stats.Float64Data
		wantErr bool
	}{
		{
			name:    "return error for non-slice type",
			raw:     struct{}{},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "return error for slice of structs that do not implement timestats.Floater",
			raw:     []struct{}{{}},
			want:    nil,
			wantErr: true,
		},
		{
			name: "transform struct that implements timestats.Floater",
			raw:  []floater{{time: 1 * time.Nanosecond}, {time: 2 * time.Second}},
			want: stats.Float64Data{1, 2000000000},
		},
		{
			name: "discard unnecessary struct fields",
			raw:  []floater{{time: 1 * time.Nanosecond, any: struct{}{}}},
			want: stats.Float64Data{1},
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			got, err := timestats.TransformIter(testcase.raw)

			if testcase.wantErr {
				if err == nil {
					t.Fatal("want error, have none")
				}
				return // Do not continue test case, transformer == nil.
			}

			if err != nil {
				t.Fatalf("want nil error, got %v", err)
			}

			if !reflect.DeepEqual(got, testcase.want) {
				t.Errorf("incorrect transform: want %v, got %v", testcase.want, got)
			}
		})
	}
}

type floatSlicer []struct {
	time time.Duration
	any  interface{}
}

func (f floatSlicer) FloatSlice() []float64 {
	floats := make([]float64, len(f))
	for i, v := range f {
		floats[i] = float64(v.time)
	}
	return floats
}

func TestTransform(t *testing.T) {
	for _, testcase := range []struct {
		name string
		raw  timestats.FloatSlicer
		want stats.Float64Data
	}{
		{
			name: "load type that implements timestats.Transform",
			raw:  floatSlicer{{time: 1 * time.Nanosecond}, {time: 2 * time.Second}},
			want: stats.Float64Data{1, 2000000000},
		},
		{
			name: "discard unnecessary struct fields",
			raw:  floatSlicer{{time: 1 * time.Nanosecond, any: struct{}{}}},
			want: stats.Float64Data{1},
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			got := timestats.Transform(testcase.raw)

			if !reflect.DeepEqual(got, testcase.want) {
				t.Errorf("incorrect transform: want %v, got %v", testcase.want, got)
			}
		})
	}
}

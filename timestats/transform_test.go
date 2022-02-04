package timestats_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/montanaflynn/stats"

	"github.com/benchttp/worker/timestats"
)

type fake struct {
	time time.Duration
	any  interface{} //nolint
}

func (e fake) Float() float64 {
	return float64(e.time)
}

func TestNewTransformer(t *testing.T) {
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
			name:    "return error for slice of structs that do not implement timestats.fake",
			raw:     []struct{}{{}},
			want:    nil,
			wantErr: true,
		},
		{
			name: "transform struct that implements timestats.fake",
			raw:  []fake{{time: 1 * time.Nanosecond}, {time: 2 * time.Second}},
			want: stats.Float64Data{1, 2000000000},
		},
		{
			name: "discard unnecessary struct fields",
			raw:  []fake{{time: 1 * time.Nanosecond, any: struct{}{}}},
			want: stats.Float64Data{1},
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			transformer, err := timestats.NewTransformer(testcase.raw)

			if testcase.wantErr {
				if err == nil {
					t.Fatal("want error, have none")
				}
				return // Do not continue test case, transformer == nil.
			}
			if err != nil {
				t.Fatalf("want nil error, got %v", err)
			}

			got := transformer.Floats
			if !reflect.DeepEqual(got, testcase.want) {
				t.Errorf("incorrect transform: want %v, got %v", testcase.want, got)
			}
		})
	}
}

type fakes []struct {
	time time.Duration
	any  interface{}
}

func (f fakes) FloatSlice() []float64 {
	floats := make([]float64, len(f))
	for i, v := range f {
		floats[i] = float64(v.time)
	}
	return floats
}

func TestSimpleTransform(t *testing.T) {
	for _, testcase := range []struct {
		name string
		raw  timestats.SimpleTransformer
		want stats.Float64Data
	}{
		{
			name: "load type that implements timestats.Transform",
			raw:  fakes{{time: 1 * time.Nanosecond}, {time: 2 * time.Second}},
			want: stats.Float64Data{1, 2000000000},
		},
		{
			name: "discard unnecessary struct fields",
			raw:  fakes{{time: 1 * time.Nanosecond, any: struct{}{}}},
			want: stats.Float64Data{1},
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			transformer := timestats.TransformSimply(testcase.raw)

			got := transformer.Floats
			if !reflect.DeepEqual(got, testcase.want) {
				t.Errorf("incorrect transform: want %v, got %v", testcase.want, got)
			}
		})
	}
}

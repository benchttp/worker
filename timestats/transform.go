package timestats

import (
	"fmt"
	"reflect"

	"github.com/montanaflynn/stats"
)

// Floater is the interface to get a float64 value from a struct.
type Floater interface {
	Float() float64
}

// TransformIter returns a data set ready for computation.
// The data set is built with floats retrieved from x.
// x must be a slice of structs implementing timestats.Floater.
func TransformIter(x interface{}) (stats.Float64Data, error) {
	s := reflect.ValueOf(x)
	if s.Kind() != reflect.Slice {
		return nil, fmt.Errorf("received a non-slice type: %s", reflect.TypeOf(x))
	}

	floats := make(stats.Float64Data, s.Len())

	for i := 0; i < s.Len(); i++ {
		v := s.Index(i)
		f, ok := v.Interface().(Floater)
		if !ok {
			return nil, fmt.Errorf("%s %w", v.Type(), ErrInterface)
		}
		floats[i] = f.Float()
	}
	return floats, nil
}

// FloatSlicer is the interface to get a slice of float64 from a type.
type FloatSlicer interface {
	FloatSlice() []float64
}

// Transform returns a data set ready for computation.
// The data set is extracted from the given FloatSlicer.
func Transform(in FloatSlicer) stats.Float64Data {
	return in.FloatSlice()
}

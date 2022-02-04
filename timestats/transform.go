package timestats

import (
	"fmt"
	"reflect"

	"github.com/montanaflynn/stats"
)

// Floater is the interface to get a float64 value from any struct.
type Floater interface {
	Float() float64
}

type Transformer struct {
	Floats stats.Float64Data
}

// NewTransformer returns a transformer loaded with floats retrieved from x.
// x is a slice of structs implementing timestats.Floater. Returns an error
// if x is not a slice or if its elements do not implement timestats.Floater.
func NewTransformer(x interface{}) (*Transformer, error) {
	s := reflect.ValueOf(x)
	if s.Kind() != reflect.Slice {
		return nil, fmt.Errorf("received a non-slice type: %s", reflect.TypeOf(x))
	}

	floats := make(stats.Float64Data, s.Len())

	for i := 0; i < s.Len(); i++ {
		v := s.Index(i)
		f, ok := v.Interface().(Floater)
		if !ok {
			return nil, fmt.Errorf("%s must implement timestats.Floater", v.Type())
		}
		floats[i] = f.Float()
	}
	return &Transformer{Floats: floats}, nil
}

// SimpleTransformer is the interface to get a slice of float64 from a type.
type SimpleTransformer interface {
	FloatSlice() []float64
}

// TransformSimply returns a transformer loaded with floats
// retrieved from the given SimpleTransformer.
func TransformSimply(in SimpleTransformer) *Transformer {
	f := in.FloatSlice()
	return &Transformer{Floats: f}
}

package timestats

import (
	"github.com/montanaflynn/stats"
)

// FloatSlicer is the interface to get a slice of float64 from a type.
type FloatSlicer interface {
	FloatSlice() []float64
}

// Transform returns a data set ready for computation.
// The data set is extracted from the given FloatSlicer.
func Transform(in FloatSlicer) stats.Float64Data {
	return in.FloatSlice()
}

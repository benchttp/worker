package timestats

import "errors"

var (
	// ErrEmptySlice is returned when an empty slice is passed to Transform.
	ErrEmptySlice = errors.New("attempted to transform empty slice")
	// ErrCompute is returned when Compute fails to compute a stat.
	ErrCompute = errors.New("failed to compute stat")
	// ErrInterface is returned when a struct that does not implement
	// Floater is passed to TransformIter.
	ErrInterface = errors.New("must implement timestats.Floater")
)

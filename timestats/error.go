package timestats

import "errors"

var (
	// ErrEmptySlice is returned when an empty slice is passed to Transform.
	ErrEmptySlice = errors.New("attempted to transform empty slice")
	// ErrCompute is returned when Compute fails to compute a stat.
	ErrCompute = errors.New("failed to compute stat")
)

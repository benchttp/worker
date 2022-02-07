package timestats

import (
	"errors"
	"fmt"
	"strings"
)

var (
	// ErrEmptySlice is returned when an empty slice is passed to Transform.
	ErrEmptySlice = errors.New("attempted to transform empty slice")
	// ErrInterface is returned when a struct that does not implement
	// Floater is passed to TransformIter.
	ErrInterface = errors.New("must implement timestats.Floater")
)

// ErrCompute is returned when Compute fails to compute a stat.
type ErrCompute struct {
	Errors []string
}

func (e *ErrCompute) Error() string {
	return fmt.Sprintf("failed to compute stat:\n  %s", strings.Join(e.Errors, "\n  "))
}

package timestats

import (
	"errors"
	"fmt"
	"strings"
)

var (
	// ErrEmptySlice is returned when working on an empty slice.
	ErrEmptySlice = errors.New("input slice is empty")
	// ErrInterface is returned when a struct that does not implement
	// Floater is passed to TransformIter.
	ErrInterface = errors.New("must implement timestats.Floater")
)

// ComputeError is returned when Compute fails to compute a stat.
type ComputeError struct {
	Errors []string
}

func (e *ComputeError) Error() string {
	return fmt.Sprintf("failed to compute stat:\n  %s", strings.Join(e.Errors, "\n  "))
}

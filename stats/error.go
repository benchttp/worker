package stats

import (
	"errors"
	"fmt"
	"strings"
)

// ErrEmptySlice is returned when working on an empty slice.
var ErrEmptySlice = errors.New("input slice is empty")

// ComputeError is returned when Compute fails to compute a stat.
type ComputeError struct {
	Errors []string
}

func (e *ComputeError) Error() string {
	return fmt.Sprintf("failed to compute stat:\n  %s", strings.Join(e.Errors, "\n  "))
}

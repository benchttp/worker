package timestats

import (
	"fmt"
)

// ordinal return x ordinal format.
//	ordinal(3) == "3rd"
func ordinal(x int) string {
	suffix := "th"
	switch x % 10 {
	case 1:
		if x%100 != 11 {
			suffix = "st"
		}
	case 2:
		if x%100 != 12 {
			suffix = "nd"
		}
	case 3:
		if x%100 != 13 {
			suffix = "rd"
		}
	}
	return fmt.Sprintf("%d%s", x, suffix)
}

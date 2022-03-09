package stats

import "fmt"

type StatusDistribution struct {
	Status1xx int
	Status2xx int
	Status3xx int
	Status4xx int
	Status5xx int
}

func ComputeStatusDistribution(data []int) (StatusDistribution, error) {
	out := StatusDistribution{}
	errs := []string{}

	for _, code := range data {
		switch code / 100 {
		case 1:
			out.Status1xx++
		case 2:
			out.Status2xx++
		case 3:
			out.Status3xx++
		case 4:
			out.Status4xx++
		case 5:
			out.Status5xx++
		default:
			errs = append(errs, fmt.Sprintf("%d is not a valid HTTP status code", code))
		}
	}

	if len(errs) > 0 {
		return out, &ComputeError{errs}
	}

	return out, nil
}

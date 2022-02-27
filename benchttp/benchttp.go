package benchttp

// Report represents the result of a Benchttp benchmark run.
type Report struct {
	Benchmark Benchmark
}

// Benchmark represents the detailed collection of requests
// done during a Benchttp benchmark run.
type Benchmark struct {
	Records []Record
}

// Times returns the recorded execution time of the requests
// as a slice of floats.
func (b Benchmark) Times() []float64 {
	s := make([]float64, len(b.Records))
	for i, v := range b.Records {
		s[i] = v.Time
	}
	return s
}

// Record represents the summary of a HTTP response.
type Record struct {
	Time float64
}

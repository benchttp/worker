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

// Values extracts and returns the status code and the execution
// time of each recorded request as two distinct slices.
func (b Benchmark) Values() (codes []int, times []float64) {
	codes = make([]int, len(b.Records))
	times = make([]float64, len(b.Records))
	for i, v := range b.Records {
		codes[i] = v.Code
		times[i] = v.Time
	}
	return
}

// Record represents the summary of a HTTP response.
type Record struct {
	Time float64
	Code int
}

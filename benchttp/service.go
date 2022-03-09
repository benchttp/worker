package benchttp

import "github.com/benchttp/worker/stats"

// InsertService defines the interface to implement by a
// data layer facade.
type InsertionService interface {
	// Create stores stats in database.
	Insert(stats.Stats, string, string, string) error
}

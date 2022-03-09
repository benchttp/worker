package benchttp

// InsertService defines the interface to implement by a
// data layer facade.
type InsertionService interface {
	// Create stores stats in database.
	Insert(Stats, string, string, string) error
}

package benchttp

// StatsService defines the interface to implement by a
// service facade inside the application.
type StatsService interface {
	// Create stores stats in database.
	Create() error
}

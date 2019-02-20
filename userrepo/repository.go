package userrepo

// Repository is a userrepo repository
type Repository interface {
	// Initialize prepares data source
	Initialize() error

	// Close closes data source and clean all resources
	Close() error

	// Validate checks if given username registered
	// and given password is correct.
	Validate(username, password string) (bool, error)
}

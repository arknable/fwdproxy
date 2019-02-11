package user

// Repository is a user repository
type Repository interface {
	// Initialize prepares data source
	Initialize() error

	// IsValid checks if given username registered
	// and given password is correct.
	IsValid(username, password string) bool

	// Random returns random username and its password
	Random() (string, string)
}

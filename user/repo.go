package user

// Repository is a user repository
type Repository interface {
	// IsValid checks if given username registered
	// and given password is correct.
	IsValid(username, password string) bool

	// Random returns random username and its password
	Random() (string, string)
}

// Repo returns instance of Repository
func Repo() Repository {
	repo := new(DummyRepository)
	repo.init()
	return repo
}

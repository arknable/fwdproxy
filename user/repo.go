package user

var (
	// BuiltInUsername is username that inserted the first time
	// database initialized.
	BuiltInUsername = "admin"

	// BuiltInUserPwd is BuiltInUsername's password.
	BuiltInUserPwd = "4dm1n"

	// Instance of Repository
	repo Repository
)

// Initialize performs initialization
func Initialize() (Repository, error) {
	repo = new(BoltRepository)
	if err := repo.Initialize(); err != nil {
		return nil, err
	}
	return repo, nil
}

// Repo returns instance of Repository
func Repo() Repository {
	return repo
}

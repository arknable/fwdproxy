package user

import "log"

var (
	// BuiltInUsername is username that inserted the first time
	// database initialized.
	BuiltInUsername = "admin"

	// BuiltInUserPwd is BuiltInUsername's password.
	BuiltInUserPwd = "4dm1n"

	// Instance of Repository
	repo Repository
)

func init() {
	repo = new(BoltRepository)
	if err := repo.Initialize(); err != nil {
		log.Fatal(err)
	}
}

// Repo returns instance of Repository
func Repo() Repository {
	return repo
}

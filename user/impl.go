package user

// Repo returns instance of Repository
func Repo() Repository {
	repo := new(DummyRepository)
	repo.Initialize()
	return repo
}

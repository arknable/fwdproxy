package user

import (
	"math/rand"
)

// DummyRepository is a dummy repo that keeps
// pre-generated users in a map
type DummyRepository struct {
	// Collection of predefined allowed users
	users map[string]string
}

// Initialize implements Repository.Initialize
func (r *DummyRepository) Initialize() error {
	r.users = map[string]string{
		"kirk":   "kirkpassword",
		"spock":  "spockpassword",
		"bones":  "bonespassword",
		"chekov": "chekovpassword",
		"sulu":   "sulupassword",
	}
	return nil
}

// IsValid implements Repository.IsValid
func (r *DummyRepository) IsValid(username, password string) bool {
	pwd, ok := r.users[username]
	if !ok {
		return false
	}
	return (password == pwd)
}

// Random implements Repository.Random
func (r *DummyRepository) Random() (string, string) {
	length := len(r.users)
	usernames := make([]string, 0, length)
	for u := range r.users {
		usernames = append(usernames, u)
	}
	var idx int
	var username string
	var password string
	ok := false

	for !ok {
		idx = rand.Intn(length - 1)
		username = usernames[idx]
		password, ok = r.users[username]
	}

	return username, password
}

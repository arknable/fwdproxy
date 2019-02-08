package config

var (
	// Collection of predefined allowed users
	users = map[string]string{
		"kirk":   "kirkpassword",
		"spock":  "spockpassword",
		"bones":  "bonespassword",
		"chekov": "chekovpassword",
		"sulu":   "sulupassword",
	}
)

// AuthIsValid checks whether given username is allowed
func AuthIsValid(username, password string) bool {
	pwd, ok := users[username]
	if !ok {
		return false
	}
	return (password == pwd)
}

// Users returns array of usernames
func Users() []string {
	usernames := make([]string, 0, len(users))
	for u := range users {
		usernames = append(usernames, u)
	}
	return usernames
}

// Password returns password of given username
func Password(username string) string {
	pwd, ok := users[username]
	if !ok {
		return ""
	}
	return pwd
}

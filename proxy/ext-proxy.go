package proxy

// External is an external proxy used to process request
type External struct {
	// Address is the listening address, its value should not
	// contains scheme, just host and port.
	Address string

	// Username is authorization's username
	Username string

	// Password is authorization's password
	Password string
}

package config

const (
	// HTTPPort is server port for HTTP
	HTTPPort = "8000"

	// TLSPort is server port for HTTPS
	TLSPort = "9000"
)

var (
	// TLSAllowedHost is host name to be allowed by ACME manager
	TLSAllowedHost = "localhost"
)

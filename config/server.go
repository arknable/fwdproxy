package config

import (
	"path"

	homedir "github.com/mitchellh/go-homedir"
)

const (
	// HTTPPort is server port for HTTP
	HTTPPort = "8000"

	// TLSPort is server port for HTTPS
	TLSPort = "9000"
)

var (
	// TLSAllowedHost is host name to be allowed by ACME manager
	TLSAllowedHost = "localhost"

	// CertCacheDir is certificate caching root folder
	CertCacheDir string

	// CertPath is path to certificate file
	CertPath string

	// KeyPath is path to private key file
	KeyPath string
)

func init() {
	CertCacheDir = "Certificates"
	home, err := homedir.Dir()
	if err == nil {
		CertCacheDir = path.Join(home, CertCacheDir)
	}
	CertPath = path.Join(CertCacheDir, "cert.pem")
	KeyPath = path.Join(CertCacheDir, "key.pem")
}

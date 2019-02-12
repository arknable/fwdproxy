package config

import (
	"path"

	homedir "github.com/mitchellh/go-homedir"
)

var (
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
	CertPath = path.Join(CertCacheDir, "pub.cert")
	KeyPath = path.Join(CertCacheDir, "key.cert")
}

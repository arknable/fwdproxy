package config

import (
	"log"
	"os"
	"path"

	homedir "github.com/mitchellh/go-homedir"
)

var (
	// App's home folder
	homePath string
)

func init() {
	folderName := ".fwdproxy"
	homePath, err := homedir.Dir()
	if err != nil {
		workPath, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		homePath = path.Join(workPath, folderName)
		return
	}

	homePath = path.Join(homePath, folderName)
}

// HomePath returns path to app's home folder
func HomePath() string {
	return homePath
}

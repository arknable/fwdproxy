package env

import (
	"os"
	"path"

	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
)

// App's home folder
var homePath string

// Initialize performs initialization
func Initialize() {
	folderName := ".fwdproxy"
	homePath, err := homedir.Dir()
	if err != nil {
		log.WithError(err).Error("Failed to find user's Home path.")

		workPath, err := os.Getwd()
		if err != nil {
			homePath = folderName
			log.WithError(err).Error("Unknown working path, using active path.")
		} else {
			homePath = path.Join(workPath, folderName)
		}
	} else {
		homePath = path.Join(homePath, folderName)
	}
	log.WithField("path", homePath).Info("Using path as Home")
}

// HomePath returns path to app's home folder
func HomePath() string {
	return homePath
}

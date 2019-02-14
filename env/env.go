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
func Initialize() error {
	folderName := ".fwdproxy"
	hpath, err := homedir.Dir()
	if err != nil {
		log.WithError(err).Error("Failed to find user's Home path.")

		workPath, err := os.Getwd()
		if err != nil {
			hpath = folderName
			log.WithError(err).Error("Unknown working path, using active path.")
		} else {
			hpath = path.Join(workPath, folderName)
		}
	} else {
		hpath = path.Join(hpath, folderName)
	}
	homePath = hpath
	_, err = os.Stat(homePath)
	if os.IsNotExist(err) {
		if err = os.MkdirAll(homePath, os.ModePerm); err != nil {
			return err
		}
	}

	log.WithField("path", homePath).Info("Using path as Home")
	return nil
}

// HomePath returns path to app's home folder
func HomePath() string {
	return homePath
}

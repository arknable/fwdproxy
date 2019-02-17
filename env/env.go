package env

import (
	"os"
	"path"

	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
)

var (
	// Users' home path
	homePath string

	// App's home path
	appHomePath string
)

// Initialize performs initialization
func Initialize() error {
	folderName := ".fwdproxy"
	hpath, err := homedir.Dir()
	if err != nil {
		log.WithError(err).Error("Failed to find user's Home path.")

		workPath, err := os.Getwd()
		if err != nil {
			appHomePath = folderName
			log.WithError(err).Error("Unknown working path, using active path.")
		} else {
			appHomePath = path.Join(workPath, folderName)
		}
	} else {
		homePath = hpath
		appHomePath = path.Join(hpath, folderName)
	}
	_, err = os.Stat(appHomePath)
	if os.IsNotExist(err) {
		if err = os.MkdirAll(appHomePath, os.ModePerm); err != nil {
			return err
		}
	}

	log.WithField("path", appHomePath).Info("Using path as Home")
	return nil
}

// HomePath returns path to app's home folder
func HomePath() string {
	return appHomePath
}

// UserHomePath returns path to users' home folder
func UserHomePath() string {
	return homePath
}

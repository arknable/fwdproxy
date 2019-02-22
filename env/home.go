package env

import (
	"os"
	"path"

	"github.com/mitchellh/go-homedir"
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
		workPath, err := os.Getwd()
		if err != nil {
			appHomePath = folderName
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
	return nil
}

// AppHomePath returns path to app's home folder
func AppHomePath() string {
	return appHomePath
}

// HomePath returns path to users' home folder
func HomePath() string {
	return homePath
}

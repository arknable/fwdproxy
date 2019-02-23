package env

import (
	"path"
	"testing"

	"github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/assert"
)

func TestInitHome(t *testing.T) {
	hpath, err := homedir.Dir()
	assert.Nil(t, err)
	appPath := path.Join(hpath, ".fwdproxy")

	err = initHome()
	assert.Nil(t, err)
	assert.Equal(t, hpath, HomePath())
	assert.Equal(t, appPath, AppHomePath())
}

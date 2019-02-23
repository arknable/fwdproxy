package env

import (
	"os"
	"path"
	"testing"

	"github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/assert"
)

func TestInitConfig(t *testing.T) {
	hpath, err := homedir.Dir()
	assert.Nil(t, err)
	hpath = path.Join(hpath, ".fwdproxy")

	_, err = os.Stat(hpath)
	if err != nil {
		if os.IsExist(err) {
			if err := os.RemoveAll(hpath); err != nil {
				t.Fatal(err)
			}
			if err := os.MkdirAll(hpath, os.ModePerm); err != nil {
				t.Fatal(err)
			}
		} else {
			t.Fatal(err)
		}
	}

	err = initConfig()
	assert.Nil(t, err)
	conf := Configuration()
	assert.NotNil(t, conf)
	assert.NotNil(t, conf.ExtProxy)
	defConf := defaultConfig()
	assert.Equal(t, conf.Port, defConf.Port)
	assert.Equal(t, conf.ExtProxy.Address, defConf.ExtProxy.Address)
	assert.Equal(t, conf.ExtProxy.Port, defConf.ExtProxy.Port)
	assert.Equal(t, conf.ExtProxy.Username, defConf.ExtProxy.Username)
	assert.Equal(t, conf.ExtProxy.Password, defConf.ExtProxy.Password)
}

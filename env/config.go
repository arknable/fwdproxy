package env

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
)

var configInstance *ServerConfiguration

// Configuration returns configuration from config file
func Configuration() *ServerConfiguration {
	return configInstance
}

// ExtProxyConfig is configuration for external proxy
type ExtProxyConfig struct {
	// Address is listening address
	Address string `json:"address"`

	// Port is listening port
	Port string `json:"port"`

	// Username is authorization username
	Username string `json:"username"`

	// Password is authorization password
	Password string `json:"password"`
}

// ServerConfiguration is configuration for proxy server
type ServerConfiguration struct {
	// Port is server port
	Port string `json:"port"`

	// ExtProxy is external proxy configuration
	ExtProxy *ExtProxyConfig `json:"proxy"`
}

// Creates default config
func createNewConfig(configPath string) error {
	file, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	conf := &ServerConfiguration{
		Port: "8000",
		ExtProxy: &ExtProxyConfig{ // By default, server assumes Tinyproxy
			Address:  "127.0.0.1",
			Port:     "8888",
			Username: "test",
			Password: "testpassword",
		},
	}
	content, err := json.MarshalIndent(conf, "", "\t")
	if err != nil {
		return err
	}
	if _, err = file.Write(content); err != nil {
		return err
	}
	return nil
}

// Initializes configuration
func initConfig() error {
	configPath := path.Join(AppHomePath(), "config.conf")
	_, err := os.Stat(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			if err = createNewConfig(configPath); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	file, err := os.Open(configPath)
	if err != nil {
		return err
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	var config ServerConfiguration
	if err = json.Unmarshal(content, &config); err != nil {
		return err
	}
	configInstance = &config

	return nil
}

package server

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"sync"

	log "github.com/sirupsen/logrus"
)

var (
	// ProxyAddress is address of external proxy server
	ProxyAddress = "http://127.0.0.1:8888"

	// ProxyUsername is username to connect to proxy
	ProxyUsername = "test"

	// ProxyPassword is password of ProxyUsername
	ProxyPassword = "testpassword"

	// Instance of transport
	trsInstance *http.Transport
	once        sync.Once
)

// Transport returns instance of http transport
func Transport() *http.Transport {
	once.Do(func() {
		instance, err := createTransport()
		if err != nil {
			log.Fatal(err)
		}
		trsInstance = instance
	})
	return trsInstance
}

// Creates instance of transport, this should be called once.
func createTransport() (*http.Transport, error) {
	proxyURL, err := url.Parse(ProxyAddress)
	if err != nil {
		return nil, err
	}
	proxyURL.User = url.UserPassword(ProxyUsername, ProxyPassword)
	t := &http.Transport{
		Proxy:        http.ProxyURL(proxyURL),
		TLSNextProto: make(map[string]func(authority string, c *tls.Conn) http.RoundTripper),
	}
	return t, nil
}

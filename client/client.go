package client

import (
	"crypto/tls"
	"net/http"
	"net/url"

	"github.com/arknable/upwork-test-proxy/config"
)

// New creates HTTP/HTTPS client given its proxy address
func New(proxyAddr string) (*http.Client, error) {
	proxyURL, err := url.Parse(proxyAddr)
	if err != nil {
		return nil, err
	}
	proxyURL.User = url.UserPassword(config.ProxyUsername, config.ProxyPassword)
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: !config.IsProduction,
			},
		},
	}

	return client, nil
}

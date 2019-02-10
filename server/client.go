package server

import (
	"crypto/tls"
	"net/http"
	"net/url"

	"github.com/arknable/fwdproxy/config"
)

// NewClient creates HTTP/HTTPS client given its proxy address
func NewClient(proxyAddr string) (*http.Client, error) {
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

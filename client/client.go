package client

import (
	"net/http"
	"net/url"
)

// NewProxied creates HTTP client that uses proxy
func NewProxied(proxyAddr string) (*http.Client, error) {
	proxyURL, err := url.Parse(proxyAddr)
	if err != nil {
		return nil, err
	}

	return &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
	}, nil
}

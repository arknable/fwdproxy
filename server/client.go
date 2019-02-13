package server

import (
	"crypto/tls"
	"net/http"

	"github.com/arknable/fwdproxy/config"
)

// NewClient creates HTTP/HTTPS client given its proxy address.
// Note: Transport.Proxy only support HTTP and SOCKS5 so we don't use it anymore:
// https://go-review.googlesource.com/c/go/+/66010/
func NewClient(proxyAddr string) (*http.Client, error) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: !config.IsProduction,
			},
		},
	}

	return client, nil
}

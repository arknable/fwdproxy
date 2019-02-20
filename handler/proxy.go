package handler

import (
	"fmt"
	"net/http"
	"net/url"
)

var (
	// ProxyAddress is external proxy address.
	// ProxyAddress shouldn't contains scheme, just host and port.
	ProxyAddress = "127.0.0.1:8888"

	// ProxyUsername is username of external proxy
	ProxyUsername = "test"

	// ProxyPassword is password of external proxy
	ProxyPassword = "testpassword"

	// HTTP transport
	pTransport *http.Transport
)

// Initialize performs initialization
func Initialize() error {
	proxyURL, err := url.Parse(fmt.Sprintf("http://%s", ProxyAddress))
	if err != nil {
		return err
	}
	proxyURL.User = url.UserPassword(ProxyUsername, ProxyPassword)

	pTransport = &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}
	return nil
}
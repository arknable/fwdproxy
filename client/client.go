package client

import (
	"crypto/tls"
	"net/http"
	"net/url"

	"github.com/arknable/upwork-test-proxy/config"
)

// NewProxied creates HTTP client that uses proxy
func NewProxied(proxyAddr string, isSecured bool) (*http.Client, error) {
	proxyURL, err := url.Parse(proxyAddr)
	if err != nil {
		return nil, err
	}
	proxyURL.User = url.UserPassword(config.ProxyUsername, config.ProxyPassword)

	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}
	if isSecured {
		// cert, err := ioutil.ReadFile(config.CertPath)
		// if err != nil {
		// 	return nil, err
		// }
		// certPool := x509.NewCertPool()
		// certPool.AppendCertsFromPEM(cert)
		// certPair, err := tls.LoadX509KeyPair(config.CertPath, config.KeyPath)
		// if err != nil {
		// 	return nil, err
		// }

		transport.TLSClientConfig = &tls.Config{
			// RootCAs:            certPool,
			// Certificates:       []tls.Certificate{certPair},
			InsecureSkipVerify: !config.IsProduction,
		}

	}

	return &http.Client{
		Transport: transport,
	}, nil
}

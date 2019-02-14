package server

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/arknable/fwdproxy/config"
)

// NewClient creates HTTP/HTTPS client given its proxy address.
// Note: Transport.Proxy only support HTTP and SOCKS5 so we don't use it anymore:
// https://go-review.googlesource.com/c/go/+/66010/
func NewClient(proxyAddr string) (*http.Client, error) {
	proxyURL, err := url.Parse(proxyAddr)
	if err != nil {
		return nil, err
	}
	proxyURL.User = url.UserPassword(config.ProxyUsername, config.ProxyPassword)

	cert, err := tls.LoadX509KeyPair(config.CertPath, config.KeyPath)
	if err != nil {
		return nil, err
	}
	caCert, err := ioutil.ReadFile(config.CACertPath)
	if err != nil {
		return nil, err
	}
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(caCert)

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
			TLSClientConfig: &tls.Config{
				Certificates:       []tls.Certificate{cert},
				RootCAs:            certPool,
				InsecureSkipVerify: !config.IsProduction,
			},
		},
	}

	return client, nil
}

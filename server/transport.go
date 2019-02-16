package server

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"sync"

	"github.com/arknable/fwdproxy/config"
	log "github.com/sirupsen/logrus"
)

var (
	// ProxyAddress is external proxy address for HTTP
	ProxyAddress = "http://127.0.0.1:8888"

	// TLSProxyAddress is external proxy address for HTTPS
	TLSProxyAddress = "http://127.0.0.1:9000"

	// ProxyUsername is username to connect to proxy
	ProxyUsername = "test"

	// ProxyPassword is password of ProxyUsername
	ProxyPassword = "testpassword"

	// HTTP transport
	transport *http.Transport
	once      sync.Once

	// TLS transport
	tlsTransport *http.Transport
	tlsOnce      sync.Once
)

// Transport returns instance of http transport
func Transport() *http.Transport {
	once.Do(func() {
		instance, err := createTransport(ProxyAddress)
		if err != nil {
			log.Fatal(err)
		}
		transport = instance
	})
	return transport
}

// TLSTransport returns instance of TLS transport
func TLSTransport() *http.Transport {
	once.Do(func() {
		cert, err := tls.LoadX509KeyPair(config.CertPath, config.KeyPath)
		if err != nil {
			log.Fatal(err)
		}

		var certPool *x509.CertPool
		if config.IsProduction {
			certPool, err = x509.SystemCertPool()
			if err != nil {
				log.Fatal(err)
			}
		} else {
			caCert, err := ioutil.ReadFile(path.Join(config.CertCacheDir, "ca.pem"))
			if err != nil {
				log.Fatal(err)
			}
			certPool = x509.NewCertPool()
			certPool.AppendCertsFromPEM(caCert)
		}

		config := &tls.Config{
			Certificates:       []tls.Certificate{cert},
			RootCAs:            certPool,
			InsecureSkipVerify: !config.IsProduction,
		}
		config.BuildNameToCertificate()

		tlsTransport = &http.Transport{
			TLSClientConfig: config,
			TLSNextProto:    make(map[string]func(authority string, c *tls.Conn) http.RoundTripper),
		}
	})
	return tlsTransport
}

// Creates instance of transport given proxy address
func createTransport(proxyAddr string) (*http.Transport, error) {
	proxyURL, err := url.Parse(proxyAddr)
	if err != nil {
		return nil, err
	}
	proxyURL.User = url.UserPassword(ProxyUsername, ProxyPassword)
	t := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}
	return t, nil
}

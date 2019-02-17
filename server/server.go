package server

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/arknable/fwdproxy/env"
)

var (
	// CertificateDir is path to certificate folder
	CertificateDir string

	// CertFileName is certificate's file name
	CertFileName = "cert.pem"

	// KeyFileName is key's file name
	KeyFileName = "key.pem"

	// CACertFileName is CA's certificate file name
	CACertFileName = "ca.pem"

	// Full path to certificate file
	certPath string

	// Full path to key file
	keyPath string

	// Port is port for HTTP listening
	Port = "8000"

	// TLSProxyAddress is external proxy address for HTTPS
	TLSProxyAddress = "https://127.0.0.1:9000"

	// ProxyAddress is external proxy address for HTTP
	ProxyAddress = "http://127.0.0.1:8888"

	// ProxyUsername is username to connect to proxy
	ProxyUsername = "test"

	// ProxyPassword is password of ProxyUsername
	ProxyPassword = "testpassword"

	// HTTP proxy URL
	proxyURL *url.URL

	// HTTPS proxy URL
	tlsProxyURL *url.URL

	// HTTP transport
	transport *http.Transport

	// HTTP server
	server *http.Server
)

// Initialize performs initialization
func Initialize(handler http.Handler) error {
	CertificateDir = path.Join(env.UserHomePath(), "Certificates")
	certPath := path.Join(CertificateDir, CertFileName)
	keyPath := path.Join(CertificateDir, KeyFileName)

	caCert, err := ioutil.ReadFile(path.Join(CertificateDir, CACertFileName))
	if err != nil {
		return err
	}
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return err
	}
	certPool, err := x509.SystemCertPool()
	if err != nil {
		return err
	}
	certPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            certPool,
		InsecureSkipVerify: !env.IsProduction,
	}
	tlsConfig.BuildNameToCertificate()

	transport = &http.Transport{
		TLSClientConfig: tlsConfig,
		TLSNextProto:    make(map[string]func(authority string, c *tls.Conn) http.RoundTripper),
	}

	pURL, err := url.Parse(ProxyAddress)
	if err != nil {
		return err
	}
	pURL.User = url.UserPassword(ProxyUsername, ProxyPassword)
	proxyURL = pURL

	pURL, err = url.Parse(TLSProxyAddress)
	if err != nil {
		return err
	}
	pURL.User = url.UserPassword(ProxyUsername, ProxyPassword)
	tlsProxyURL = pURL

	server = &http.Server{
		Addr:         fmt.Sprintf("127.0.0.1:%s", Port),
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 20 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      handler,
		TLSConfig:    tlsConfig,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}

	return nil
}

// Start starts the server
func Start() error {
	return server.ListenAndServeTLS(certPath, keyPath)
}

// NewClient creates new HTTP client
func NewClient(useTLS bool) *http.Client {
	if useTLS {
		transport.Proxy = http.ProxyURL(tlsProxyURL)
	} else {
		transport.Proxy = http.ProxyURL(proxyURL)
	}

	return &http.Client{
		Transport: transport,
	}
}

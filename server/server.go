package server

import (
	"crypto/tls"
	"fmt"
	"github.com/arknable/fwdproxy/env"
	"log"
	"net"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

var (
	// HttpPort is port for HTTP listener
	HttpPort = "8000"

	// TlsPort is port for HTTPS listener
	TlsPort = "9000"

	// CertificatePath is path to certificate folder
	CertificatePath string

	// Path to certificate file
	certPath string

	// Path to certificate's key file
	certKeyPath string

	// HTTP transport
	transport *http.Transport

	// HTTP server
	httpServer *http.Server

	// HTTPS server
	tlsServer *http.Server
)

// Initialize performs initialization
func Initialize(handler http.Handler) error {
	CertificatePath = strings.Trim(CertificatePath, " ")
	if len(CertificatePath) == 0 {
		CertificatePath = path.Join(env.UserHomePath(), "Certificates")
	}
	certPath = path.Join(CertificatePath, "cert.pem")
	certKeyPath = path.Join(CertificatePath, "key.pem")

	proxyURL, err := url.Parse(fmt.Sprintf("http://%s", ProxyAddress))
	if err != nil {
		return err
	}
	proxyURL.User = url.UserPassword(ProxyUsername, ProxyPassword)

	transport = &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}

	httpServer = &http.Server{
		Addr: net.JoinHostPort("", HttpPort),
		IdleTimeout: 1 * time.Minute,
		ReadTimeout: 1 * time.Minute,
		WriteTimeout: 1 * time.Minute,
		Handler: handler,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}

	tlsServer = &http.Server{
		Addr: net.JoinHostPort("", TlsPort),
		IdleTimeout: 1 * time.Minute,
		ReadTimeout: 1 * time.Minute,
		WriteTimeout: 1 * time.Minute,
		Handler: handler,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}

	return nil
}

// NewClient creates new HTTP client
func NewClient() *http.Client {
	return &http.Client {
		Transport: transport,
		Timeout: 30 * time.Second,
	}
}

// Start initiate port listening
func Start() {
	go func() {
		if err := tlsServer.ListenAndServeTLS(certPath, certKeyPath); err != nil {
			log.Fatal(err)
		}
	}()

	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

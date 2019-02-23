package proxy

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/arknable/fwdproxy/env"
)

// Server is an HTTP server
type Server struct {
	Port       string
	httpServer *http.Server

	// Transport for http request
	transport *http.Transport

	// Encoded proxy credential
	proxyAuthEncoded string
}

// New creates new server
func New() (*Server, error) {
	config := env.Configuration()
	proxyConfig := config.ExtProxy

	proxyURL, err := url.Parse(fmt.Sprintf("http://%s", net.JoinHostPort(proxyConfig.Address, proxyConfig.Port)))
	if err != nil {
		return nil, err
	}
	proxyURL.User = url.UserPassword(proxyConfig.Username, proxyConfig.Password)
	srv := &Server{
		Port: config.Port,
		transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
	}
	srv.httpServer = &http.Server{
		Addr:         net.JoinHostPort("", srv.Port),
		IdleTimeout:  1 * time.Minute,
		ReadTimeout:  1 * time.Minute,
		WriteTimeout: 1 * time.Minute,
		Handler:      srv,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}

	proxyCred := fmt.Sprintf("%s:%s", proxyConfig.Username, proxyConfig.Password)
	srv.proxyAuthEncoded = base64.StdEncoding.EncodeToString([]byte(proxyCred))

	return srv, nil
}

// Start starts HTTP server
func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

// NewClient creates new HTTP client
func (s *Server) NewClient() *http.Client {
	return &http.Client{
		Transport: s.transport,
		Timeout:   30 * time.Second,
	}
}

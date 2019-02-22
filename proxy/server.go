package proxy

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"
)

// Server is an HTTP server
type Server struct {
	Port       string
	proxy      *External
	transport  *http.Transport
	httpServer *http.Server
}

// New creates new server
func New() (*Server, error) {
	proxyConfig := &External{
		Address:  "127.0.0.1:8888",
		Username: "test",
		Password: "testpassword",
	}
	proxyURL, err := url.Parse(fmt.Sprintf("http://%s", proxyConfig.Address))
	if err != nil {
		return nil, err
	}
	proxyURL.User = url.UserPassword(proxyConfig.Username, proxyConfig.Password)
	srv := &Server{
		Port:  "8000",
		proxy: proxyConfig,
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

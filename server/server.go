package server

import (
	"fmt"
	"net/http"
	"time"
)

var (
	// HTTPPort is port for HTTP listening
	HTTPPort = "8000"

	// TLSPort is port for HTTPS listening
	TLSPort = "8001"
)

// New creates an HTTP server
func New(handler http.Handler) *http.Server {
	return &http.Server{
		Addr:         fmt.Sprintf("127.0.0.1:%s", HTTPPort),
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 20 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      handler,
	}
}

// NewTLS creates an HTTPS server
func NewTLS(handler http.Handler) *http.Server {
	srv := New(handler)
	srv.Addr = fmt.Sprintf("127.0.0.1:%s", TLSPort)
	return srv
}

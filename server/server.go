package server

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/arknable/fwdproxy/config"
)

// New creates an HTTP server
func New(handler http.Handler) *http.Server {
	return &http.Server{
		Addr:         fmt.Sprintf(":%s", config.HTTPPort),
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 20 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      handler,
	}
}

// NewTLS creates an HTTPS server
func NewTLS(handler http.Handler) *http.Server {
	srv := New(handler)
	srv.Addr = fmt.Sprintf(":%s", config.TLSPort)
	srv.TLSNextProto = make(map[string]func(*http.Server, *tls.Conn, http.Handler))
	return srv
}

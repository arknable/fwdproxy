package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/arknable/upwork-test-proxy/handler"
)

const (
	// HTTPPort is server port for HTTP
	HTTPPort = "8000"

	// TLSPort is server port for HTTPS
	TLSPort = "9000"
)

// New creates an HTTP server
func New() *http.Server {
	return &http.Server{
		Addr:         fmt.Sprintf(":%s", HTTPPort),
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 20 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      http.HandlerFunc(handler.HandleRequest),
	}
}

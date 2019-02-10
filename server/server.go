package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/arknable/fwdproxy/config"
	"golang.org/x/crypto/acme/autocert"
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
func NewTLS(handler http.Handler) (*autocert.Manager, *http.Server) {
	srv := New(handler)
	srv.Addr = fmt.Sprintf(":%s", config.TLSPort)
	if !config.IsProduction {
		return nil, srv
	}

	if _, err := os.Stat(config.CertCacheDir); os.IsNotExist(err) {
		os.MkdirAll(config.CertCacheDir, os.ModePerm)
	}
	manager := &autocert.Manager{
		Prompt: autocert.AcceptTOS,
		Cache:  autocert.DirCache(config.CertCacheDir),
		HostPolicy: func(ctx context.Context, host string) error {
			if host == config.TLSAllowedHost {
				return nil
			}
			return fmt.Errorf("AutoCert/ACME: only %s allowed", config.TLSAllowedHost)
		},
	}
	srv.TLSConfig = &tls.Config{
		GetCertificate: manager.GetCertificate,
	}
	return manager, srv
}

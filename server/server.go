package server

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"
)

var (
	// Port is port for HTTP listener
	Port = "8000"

	// HTTP server
	server *http.Server

	// Client transport
	transport *http.Transport
)

// Initialize performs initialization
func Initialize(handler http.Handler) error {
	proxyURL, err := url.Parse(fmt.Sprintf("http://%s", ProxyAddress))
	if err != nil {
		return err
	}
	proxyURL.User = url.UserPassword(ProxyUsername, ProxyPassword)
	transport = &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}
	server = &http.Server{
		Addr:         net.JoinHostPort("", Port),
		IdleTimeout:  1 * time.Minute,
		ReadTimeout:  1 * time.Minute,
		WriteTimeout: 1 * time.Minute,
		Handler:      handler,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}

	return nil
}

// NewClient creates new HTTP client
func NewClient() *http.Client {
	return &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}
}

// Start initiate port listening
func Start() {
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

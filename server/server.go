package server

import (
	"fmt"
	"net/http"
	"time"
)

// Port is server port
var Port = "8000"

// New creates an HTTP server
func New(handler http.Handler) *http.Server {
	return &http.Server{
		Addr:         fmt.Sprintf("127.0.0.1:%s", Port),
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 20 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      handler,
	}
}

// NewTLS creates an HTTPS server
// func NewTLS(handler http.Handler) *http.Server {
// 	srv := New(handler)
// 	srv.Addr = fmt.Sprintf(":%s", config.TLSPort)
// 	srv.TLSNextProto = make(map[string]func(*http.Server, *tls.Conn, http.Handler))
// 	return srv
// }

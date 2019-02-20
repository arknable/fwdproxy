package main

import (
	"crypto/tls"
	"github.com/arknable/fwdproxy/handler"
	"log"
	"net"
	"net/http"
	"time"
)

var (
	// Port is server port
	Port = "8000"
)

func main() {
	srv := &http.Server{
		Addr: net.JoinHostPort("", Port),
		IdleTimeout: 1 * time.Minute,
		ReadTimeout: 1 * time.Minute,
		WriteTimeout: 1 * time.Minute,
		Handler: http.HandlerFunc(handler.Serve),
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}

	log.Printf("Listening at %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

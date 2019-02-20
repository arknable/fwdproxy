package main

import (
	"crypto/tls"
	"github.com/arknable/fwdproxy/env"
	"github.com/arknable/fwdproxy/handler"
	"github.com/arknable/fwdproxy/userrepo"
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
	if err := env.Initialize(); err != nil {
		log.Fatal(err)
	}
	repo, err := userrepo.Initialize()
	if err != nil {
		log.Fatal(err)
	}
	defer repo.Close()

	httpServer := &http.Server{
		Addr: net.JoinHostPort("", Port),
		IdleTimeout: 1 * time.Minute,
		ReadTimeout: 1 * time.Minute,
		WriteTimeout: 1 * time.Minute,
		Handler: http.HandlerFunc(handler.Serve),
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}

	log.Printf("Listening at %s", httpServer.Addr)
	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

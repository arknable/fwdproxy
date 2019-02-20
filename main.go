package main

import (
	"github.com/arknable/fwdproxy/handler"
	"log"
	"net"
	"net/http"
)

var (
	// Port is server port
	Port = "8000"
)

func main() {
	addr := net.JoinHostPort("", Port)
	log.Printf("Listening at %s", addr)
	if err := http.ListenAndServe(addr, http.HandlerFunc(handler.Serve)); err != nil {
		log.Fatal(err)
	}
}

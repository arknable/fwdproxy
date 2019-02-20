package main

import (
	"github.com/arknable/fwdproxy/env"
	"github.com/arknable/fwdproxy/handler"
	"github.com/arknable/fwdproxy/server"
	"github.com/arknable/fwdproxy/userrepo"
	"log"
	"net/http"
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

	if err := server.Initialize(http.HandlerFunc(handler.Serve)); err != nil {
		log.Fatal(err)
	}
	server.Start()
}

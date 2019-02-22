package main

import (
	"log"
	"net/http"

	"github.com/arknable/fwdproxy/env"
	"github.com/arknable/fwdproxy/handler"
	"github.com/arknable/fwdproxy/server"
	"github.com/arknable/fwdproxy/userrepo"
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

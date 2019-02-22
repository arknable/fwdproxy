package main

import (
	"log"

	"github.com/arknable/fwdproxy/env"
	"github.com/arknable/fwdproxy/proxy"
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

	proxy, err := proxy.New()
	if err != nil {
		log.Fatal(err)
	}
	if err = proxy.Start(); err != nil {
		log.Fatal(err)
	}
}

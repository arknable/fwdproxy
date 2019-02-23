package main

import (
	"log"
	"os"
	"path"

	"github.com/arknable/fwdproxy/env"
	"github.com/arknable/fwdproxy/proxy"
	"github.com/arknable/fwdproxy/userrepo"
)

func main() {
	if err := env.Initialize(); err != nil {
		log.Fatal(err)
	}

	// Setup log output to file and stdout
	logPath := path.Join(env.AppHomePath(), "output.log")
	file, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	log.SetOutput(file)

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

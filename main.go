package main

import (
	"log"
	"os"

	"github.com/arknable/upwork-test-proxy/server"
)

func main() {
	log.SetOutput(os.Stdout)

	srv := server.New()
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

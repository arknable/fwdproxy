package main

import (
	"log"
	"net/http"
	"os"

	"github.com/arknable/upwork-test-proxy/handler"
)

func main() {
	log.SetOutput(os.Stdout)
	http.ListenAndServe(":8234", http.HandlerFunc(handler.HandleRequest))
}

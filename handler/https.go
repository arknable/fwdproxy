package handler

import (
	"log"
	"net/http"
)

// Handles HTTPS request
func serveTLS(w http.ResponseWriter, r *http.Request) {
	log.Println("serveTLS")
}

package handler

import (
	"net/http"
	"strings"
)

// Serve is parent handler that receives request from server
func Serve(w http.ResponseWriter, r *http.Request) {
	if strings.ToLower(r.URL.Scheme) == "http" {
		serveHTTP(w, r)
		return
	}
	serveTLS(w, r)
}
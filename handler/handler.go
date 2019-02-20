package handler

import "net/http"

// Serve is parent handler that receives request from server
func Serve(w http.ResponseWriter, r *http.Request) {
	serveHTTP(w, r)
}
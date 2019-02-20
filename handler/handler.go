package handler

import "net/http"

// Serve is parent handler that receives request from server
func Serve(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodConnect {
		serveTLS(w, r)
		return
	}
	serveHTTP(w, r)
}
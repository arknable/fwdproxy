package handler

import (
	"github.com/arknable/fwdproxy/server"
	"io"
	"net/http"
)

// Handles HTTP request
func serveHTTP(w http.ResponseWriter, r *http.Request) {
	if err := authenticate(r); err != nil {
		status := http.StatusInternalServerError
		if err == ErrAuthRequired {
			status = http.StatusUnauthorized
		} else if err == ErrForbidden {
			status = http.StatusForbidden
		}
		http.Error(w, err.Error(), status)
		return
	}

	request, err := http.NewRequest(r.Method, r.URL.String(), r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	copyHeader(r.Header, request.Header)
	request.Header.Del("Proxy-Authorization")

	client := server.NewClient()
	response, err := client.Do(request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	copyHeader(response.Header, w.Header())
	_, err = io.Copy(w, response.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

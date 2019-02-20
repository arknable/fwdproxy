package handler

import (
	"io"
	"net/http"
	"time"
)

// Handles HTTP request
func serveHTTP(w http.ResponseWriter, r *http.Request) {
	request, err := http.NewRequest(r.Method, r.URL.String(), r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	copyHeader(r.Header, request.Header)
	request.Header.Del("Proxy-Authorization")
	client := &http.Client { Timeout: 30 * time.Second}
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

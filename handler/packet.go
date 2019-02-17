package handler

import (
	"io"
	"net/http"
)

// Copy header from src to dest
func copyHeader(src, dest http.Header) {
	for key, val := range src {
		for _, v := range val {
			dest.Add(key, v)
		}
	}
}

// Copy response from src to dest
func copyResponse(src *http.Response, dest http.ResponseWriter) error {
	dest.WriteHeader(src.StatusCode)
	copyHeader(src.Header, dest.Header())
	_, err := io.Copy(dest, src.Body)
	return err
}

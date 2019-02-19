package handler

import (
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
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
	copyHeader(src.Header, dest.Header())
	_, err := io.Copy(dest, src.Body)
	return err
}

// Transfers connection from src to dest
func transfer(src io.ReadCloser, dest io.WriteCloser) {
	_, err := io.Copy(dest, src)
	if err != nil {
		log.WithError(err).Error("Connection transfer failed")
	}
}

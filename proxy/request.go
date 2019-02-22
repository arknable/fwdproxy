package proxy

import (
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

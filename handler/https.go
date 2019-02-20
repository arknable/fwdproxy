package handler

import (
	"fmt"
	"net/http"
)

// Handles HTTPS request
func serveTLS(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w,"serveTLS")
}

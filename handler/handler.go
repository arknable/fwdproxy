package handler

import (
	"net/http"

	"github.com/arknable/upwork-test-proxy/config"
	httphandler "github.com/arknable/upwork-test-proxy/http"
	tlshandler "github.com/arknable/upwork-test-proxy/tls"
)

// HandleRequest handles requests
func HandleRequest(res http.ResponseWriter, req *http.Request) {
	username, password, ok := req.BasicAuth()
	if !ok {
		http.Error(res, "Missing authorization", http.StatusUnauthorized)
		return
	}
	if !config.AuthIsValid(username, password) {
		http.Error(res, "Invalid authorization", http.StatusForbidden)
		return
	}

	if req.Method == http.MethodConnect {
		tlshandler.HandleRequest(res, req)
		return
	}
	httphandler.HandleRequest(res, req)
}

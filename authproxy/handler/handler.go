package handler

import (
	"net/http"

	cfg "github.com/arknable/upwork-test-proxy/authproxy/config"
	httphandler "github.com/arknable/upwork-test-proxy/authproxy/http"
	tlshandler "github.com/arknable/upwork-test-proxy/authproxy/tls"
)

// HandleRequest handles requests
func HandleRequest(res http.ResponseWriter, req *http.Request) {
	username, password, ok := req.BasicAuth()
	if !ok {
		http.Error(res, "Missing authorization", http.StatusUnauthorized)
		return
	}
	if (username != cfg.AllowedUsername) || (password != cfg.AllowedPassword) {
		http.Error(res, "Invalid authorization", http.StatusForbidden)
		return
	}

	if req.Method == http.MethodConnect {
		tlshandler.HandleRequest(res, req)
		return
	}
	httphandler.HandleRequest(res, req)
}

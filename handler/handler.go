package handler

import (
	"net/http"
	"strings"

	"github.com/arknable/upwork-test-proxy/config"
	httphandler "github.com/arknable/upwork-test-proxy/http"
)

// HandleRequest handles requests
func HandleRequest(res http.ResponseWriter, req *http.Request) {
	if req.URL.Scheme == "http" {
		username, password, ok := req.BasicAuth()
		if !ok || (len(strings.Trim(username, " ")) == 0) {
			http.Error(res, "Restricted access only", http.StatusUnauthorized)
			return
		}
		if !config.AuthIsValid(username, password) {
			http.Error(res, "You have no access to do a request", http.StatusForbidden)
			return
		}
	} else {
		req.URL.Scheme = "https"
	}

	// if req.Method == http.MethodConnect {
	// 	tlshandler.HandleRequest(res, req)
	// 	return
	// }
	httphandler.HandleRequest(res, req)
}

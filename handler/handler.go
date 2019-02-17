package handler

import (
	"errors"
	"net/http"

	"github.com/arknable/fwdproxy/user"
)

// ErrInvalidAuth is error message when supplied proxy authentication is invalid.
var ErrInvalidAuth = errors.New("Invalid authentication")

// HandleRequest handles both HTTP and HTTPS request
func HandleRequest(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodConnect {
		handleTLS(res, req)
		return
	}
	handleHTTP(res, req)
}

// Checks request whether it attach Proxy-Authorization information
// and then validate it's username & password.
func validateRequest(req *http.Request) (username, password string, err error) {
	username, password, err = proxyAuth(req)
	if err != nil {
		return
	}
	valid, err := user.Repo().Validate(username, password)
	if err != nil {
		return
	}
	if !valid {
		err = ErrInvalidAuth
		return
	}
	return
}

package handler

import (
	"errors"
	"github.com/arknable/fwdproxy/userrepo"
	"net/http"
)

var (
	// ErrAuthRequired occurs when request does not have Proxy-Authorization
	ErrAuthRequired = errors.New("Proxy authentication required")

	// ErrForbidden occurs when proxy authentication is not valid
	ErrForbidden = errors.New("Forbidden access")
)

// Checks Proxy-Authorization from request header.
func authenticate(req *http.Request) error {
	username, password, ok := req.BasicAuth()
	if !ok {
		return ErrAuthRequired
	}

	isValid, err := userrepo.Instance().Validate(username, password)
	if err != nil {
		return err
	}
	if !isValid {
		return ErrForbidden
	}

	return nil
}

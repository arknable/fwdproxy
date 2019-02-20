package handler

import (
	"encoding/base64"
	"errors"
	"github.com/arknable/fwdproxy/userrepo"
	"net/http"
	"strings"
)

var (
	// ErrAuthRequired occurs when request does not have Proxy-Authorization
	ErrAuthRequired = errors.New("Proxy authentication required")

	// ErrForbidden occurs when proxy authentication is not valid
	ErrForbidden = errors.New("Forbidden access")
)

// Checks Proxy-Authorization from request header.
func authenticate(req *http.Request) error {
	header := req.Header.Get("Proxy-Authorization")
	if !strings.HasPrefix(header,"Basic") {
		return ErrAuthRequired
	}
	header = strings.TrimPrefix(header,"Basic ")
	decoded, err := base64.StdEncoding.DecodeString(header)
	decodedString := string(decoded)
	if !strings.Contains(decodedString, ":") {
		return ErrAuthRequired
	}
	credentials := strings.Split(decodedString, ":")
	username := credentials[0]
	password := credentials[1]

	isValid, err := userrepo.Instance().Validate(username, password)
	if err != nil {
		return err
	}
	if !isValid {
		return ErrForbidden
	}

	return nil
}

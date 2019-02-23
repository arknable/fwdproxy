package proxy

import (
	"encoding/base64"
	"errors"
	"net/http"
	"strings"

	"github.com/arknable/fwdproxy/userrepo"
)

var (
	// ErrAuthRequired occurs when request does not have Proxy-Authorization
	ErrAuthRequired = errors.New("Proxy authentication required")

	// ErrForbidden occurs when proxy authentication is not valid
	ErrForbidden = errors.New("Forbidden access")
)

// Authenticated checks Proxy-Authorization from request header.
func (c *Context) Authenticated() bool {
	req := c.request
	header := req.Header.Get("Proxy-Authorization")
	if !strings.HasPrefix(header, "Basic") {
		c.ResponseError(ErrAuthRequired, http.StatusUnauthorized)
		return false
	}
	header = strings.TrimPrefix(header, "Basic ")
	decoded, err := base64.StdEncoding.DecodeString(header)
	decodedString := string(decoded)
	if !strings.Contains(decodedString, ":") {
		c.ResponseError(ErrAuthRequired, http.StatusUnauthorized)
		return false
	}
	credentials := strings.Split(decodedString, ":")
	username := credentials[0]
	password := credentials[1]

	isValid, err := userrepo.Instance().Validate(username, password)
	if err != nil {
		c.ResponseError(err, http.StatusUnauthorized)
		return false
	}
	if !isValid {
		c.ResponseError(ErrForbidden, http.StatusForbidden)
		return false
	}

	return isValid
}

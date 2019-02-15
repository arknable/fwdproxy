package handler

import (
	"encoding/base64"
	"net/http"
	"strings"
)

// Extracts Proxy-Authorization from request header.
func parseProxyAuth(req *http.Request) (string, string, error) {
	stringValue := req.Header.Get("Proxy-Authorization")
	stringValue = strings.Trim(stringValue, "Basic")
	stringValue = strings.Trim(stringValue, " ")
	data, err := base64.StdEncoding.DecodeString(stringValue)
	if err != nil {
		return "", "", err
	}
	parts := strings.Split(string(data), ":")
	username := parts[0]
	password := parts[1]

	return username, password, nil
}

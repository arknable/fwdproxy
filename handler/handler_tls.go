package handler

import (
	"net/http"

	mylog "github.com/arknable/fwdproxy/log"
	"github.com/arknable/fwdproxy/server"
	log "github.com/sirupsen/logrus"
)

// Handles TLS tunneling
func handleTLS(res http.ResponseWriter, req *http.Request) {
	req.URL.Scheme = "https" // If method is CONNECT, URL.Scheme usually cleared.

	username, password, err := validateRequest(req)
	credFields := log.Fields{
		"username": username,
		"pasword":  password,
	}
	if err != nil {
		if err == ErrInvalidAuth {
			http.Error(res, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			log.WithFields(credFields).Warning(http.StatusText(http.StatusUnauthorized))
			return
		}

		internalError(res, req, err)
		return
	}
	mylog.WithRequest(req).WithFields(credFields).Info("Authenticated")

	request, err := http.NewRequest(http.MethodGet, req.URL.String(), req.Body)
	if err != nil {
		internalError(res, req, err)
		return
	}
	copyHeader(req.Header, request.Header)
	request.Header.Del("Proxy-Authorization")
	request.Header.Del("Proxy-Connection")
	// request.Header.Set("Host", req.URL.String())

	// cred := base64.StdEncoding.EncodeToString([]byte(net.JoinHostPort(server.ProxyUsername, server.ProxyPassword)))
	// request.Header.Add("Proxy-Authorization", "Basic "+cred)

	mylog.WithRequest(request).Info("Forwarded request")

	client := server.NewClient()
	resp, err := client.Do(request)
	if err != nil {
		internalError(res, req, err)
		return
	}
	defer resp.Body.Close()

	resp, err = client.Get(req.URL.String())
	if err != nil {
		internalError(res, req, err)
		return
	}
	defer resp.Body.Close()
	mylog.WithResponse(resp).Info("Returned response")
	if err = copyResponse(resp, res); err != nil {
		internalError(res, req, err)
		return
	}
}

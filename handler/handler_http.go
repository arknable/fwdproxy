package handler

import (
	"net/http"

	mylog "github.com/arknable/fwdproxy/log"
	"github.com/arknable/fwdproxy/server"
	"github.com/arknable/fwdproxy/user"
	log "github.com/sirupsen/logrus"
)

// Handles HTTP requests
func handleHTTP(res http.ResponseWriter, req *http.Request) {
	username, password, err := proxyAuth(req)
	if err != nil {
		internalError(res, req, err)
		return
	}
	mylog.WithRequest(req).Info("Incoming HTTP request")

	credFields := log.Fields{
		"username": username,
		"pasword":  password,
	}
	valid, err := user.Repo().Validate(username, password)
	if err != nil {
		internalError(res, req, err)
		return
	}
	if !valid {
		http.Error(res, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		log.WithFields(credFields).Warning(http.StatusText(http.StatusForbidden))
		return
	}
	log.WithFields(credFields).Info("Authenticated")

	request, err := http.NewRequest(req.Method, req.URL.String(), req.Body)
	if err != nil {
		internalError(res, req, err)
		return
	}
	copyHeader(req.Header, request.Header)
	request.Header.Del("Proxy-Authorization")
	mylog.WithRequest(request).Info("Forwarded request")

	client := server.NewClient()
	resp, err := client.Do(request)
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

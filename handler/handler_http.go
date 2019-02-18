package handler

import (
	"net/http"

	mylog "github.com/arknable/fwdproxy/log"
	"github.com/arknable/fwdproxy/server"
	log "github.com/sirupsen/logrus"
)

// Handles HTTP requests
func handleHTTP(res http.ResponseWriter, req *http.Request) {
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

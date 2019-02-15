package handler

import (
	"io"
	"net/http"

	mylog "github.com/arknable/fwdproxy/log"
	"github.com/arknable/fwdproxy/server"
	"github.com/arknable/fwdproxy/user"
	log "github.com/sirupsen/logrus"
)

// Handles TLS tunneling
func handleTLS(res http.ResponseWriter, req *http.Request) {
	username, password, err := parseProxyAuth(req)
	if err != nil {
		msg := "Unauthorized request"
		http.Error(res, msg, http.StatusUnauthorized)
		mylog.WithRequest(req).Warning(msg)
		return
	}

	mylog.WithRequest(req).Info("Incoming HTTPS request")

	credFields := log.Fields{
		"username": username,
		"pasword":  password,
	}
	valid, err := user.Repo().Validate(username, password)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		log.WithFields(credFields).Error(err)
		return
	}
	if !valid {
		msg := "Forbidden request"
		http.Error(res, msg, http.StatusForbidden)
		log.WithFields(credFields).Warning(msg)
		return
	}
	log.WithFields(credFields).Info("Authenticated")

	req.URL.Scheme = "https"
	request, err := http.NewRequest(http.MethodGet, req.URL.String(), req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}

	for key, val := range req.Header {
		if key == "Proxy-Authorization" {
			continue
		}

		for _, v := range val {
			request.Header.Add(key, v)
		}
	}

	mylog.WithRequest(request).Info("Forwarded HTTPS request")

	client := &http.Client{Transport: server.TLSTransport()}
	resp, err := client.Do(request)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	defer resp.Body.Close()

	mylog.WithResponse(resp).Info("Returned HTTPS response")

	for key, val := range resp.Header {
		for _, v := range val {
			res.Header().Add(key, v)
		}
	}
	res.WriteHeader(resp.StatusCode)

	_, err = io.Copy(res, resp.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
}

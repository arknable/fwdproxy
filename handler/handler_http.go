package handler

import (
	"io"
	"net/http"

	mylog "github.com/arknable/fwdproxy/log"
	"github.com/arknable/fwdproxy/server"
	"github.com/arknable/fwdproxy/user"
	log "github.com/sirupsen/logrus"
)

// Handles HTTP requests
func handleHTTP(res http.ResponseWriter, req *http.Request) {
	authExists := true
	username, password, ok := req.BasicAuth()
	if !ok {
		// Check if user uses proxy authentication
		proxUser, proxPwd, err := parseProxyAuth(req)
		if err != nil {
			authExists = false
		} else {
			username = proxUser
			password = proxPwd
		}
	}
	if !authExists {
		msg := "Unauthorized request"
		http.Error(res, msg, http.StatusUnauthorized)
		mylog.WithRequest(req).Warning(msg)
		return
	}

	mylog.WithRequest(req).Info("Incoming request")

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

	request, err := http.NewRequest(req.Method, req.URL.String(), req.Body)
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

	mylog.WithRequest(request).Info("Forwarded request")

	client := &http.Client{Transport: server.Transport()}
	resp, err := client.Do(request)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	defer resp.Body.Close()

	mylog.WithResponse(resp).Info("Returned response")

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

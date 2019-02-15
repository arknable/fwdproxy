package handler

import (
	"net/http"

	mylog "github.com/arknable/fwdproxy/log"
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
}

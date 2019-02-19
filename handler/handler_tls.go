package handler

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
	"strings"

	mylog "github.com/arknable/fwdproxy/log"
	"github.com/arknable/fwdproxy/server"
	log "github.com/sirupsen/logrus"
)

// Handles TLS tunneling
func handleTLS(res http.ResponseWriter, req *http.Request) {
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

	proxyConn, err := net.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		internalError(res, req, err)
		return
	}
	defer proxyConn.Close()

	var request bytes.Buffer
	_, err = fmt.Fprintf(&request, "CONNECT %s HTTP/1.1\r\n", req.RequestURI)
	if err != nil {
		internalError(res, req, err)
		return
	}

	header := map[string]string{
		"Host": req.RequestURI,
	}
	proxyAuth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", server.ProxyUsername, server.ProxyPassword)))
	header["Proxy-Authorization"] = "Basic " + proxyAuth
	header["Proxy-Connection"] = "Keep-Alive"
	header["Connection"] = "Keep-Alive"
	header["User-Agent"] = req.Header.Get("User-Agent")

	for k, v := range header {
		_, err = fmt.Fprintf(&request, "%s: %s\r\n", k, v)
		if err != nil {
			internalError(res, req, err)
			return
		}
	}
	_, err = fmt.Fprint(&request, "\r\n")
	if err != nil {
		internalError(res, req, err)
		return
	}
	log.Debug(request.String())
	_, err = proxyConn.Write(request.Bytes())
	if err != nil {
		internalError(res, req, err)
		return
	}
	status, err := bufio.NewReader(proxyConn).ReadString('\n')
	if err != nil {
		internalError(res, req, err)
		return
	}

	if !strings.Contains(status, "200") {
		http.Error(res, status, http.StatusInternalServerError)
		return
	}
}
